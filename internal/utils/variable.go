package utils

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var ErrBadVariable = errors.New("bad variable")

func ParseVariable(v string) ([]string, error) {
	rg := regexp.MustCompile(`\${(?<variable>.+)}`)

	matches := rg.FindStringSubmatch(v)
	if len(matches) < 2 {
		return nil, fmt.Errorf("%w: %w", ErrBadVariable, fmt.Errorf("invalid variable '%s'", v))
	}

	variable := matches[rg.SubexpIndex("variable")]
	return strings.Split(variable, "::"), nil
}

type VariableParseFunc func(params []string) (string, error)

func ReplaceVariable(envs map[string]string, variableParseFunc VariableParseFunc) (map[string]string, error) {
	rg := regexp.MustCompile(`\$\{([^\}]*)\}?`)

	replacedEnvs := map[string]string{}
	for key, val := range envs {
		ctx, cancel := context.WithCancelCause(context.Background())

		replaced := rg.ReplaceAllStringFunc(val, func(variable string) string {
			if ctx.Err() != nil {
				return variable
			}

			params, err := ParseVariable(variable)
			if err != nil {
				cancel(err)
				return variable
			}

			parsedVariable, err := variableParseFunc(params)
			if err != nil {
				cancel(err)
				return variable
			}

			return parsedVariable
		})

		if err := context.Cause(ctx); err != nil {
			return nil, err
		}

		replacedEnvs[key] = replaced
	}

	return replacedEnvs, nil
}
