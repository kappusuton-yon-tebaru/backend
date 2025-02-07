package dockerhub

import (
	// "bytes"
	// "encoding/json"
	// "fmt"
	// "net/http"
	"os"
)

type DockerHubRepository struct {
	apiURL string
	token  string
}

func NewDockerHubRepository(apiURL, token string) *DockerHubRepository {
	return &DockerHubRepository{
		apiURL: "https://hub.docker.com/v2",
		token:  os.Getenv("DOCKERHUB_TOKEN"),
	}
}

// func (r *DockerHubRepository) GetImages() ([]string, error) {
// 	url := fmt.Sprintf("%s/repositories/%s/images", r.apiURL, "kappusuton-yon-tebaru")
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.token))

// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer res.Body.Close()

// 	var result struct {
// 		Results []struct {
// 			Name string `json:"name"`
// 		} `json:"results"`
// 	}

// 	json.NewDecoder(res.Body).Decode(&result)

// 	var images []string
// 	for _, img := range result.Results {
// 		images = append(images, img.Name)
// 	}
// 	return images, nil
	
// }

// func (r *DockerHubRepository) PushImage(imageName, tag string) error {
// 	url := fmt.Sprintf("%s/repositories/%s/images", r.apiURL, "kappusuton-yon-tebaru")
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(nil))
// 	if err != nil {
// 		return err
// 	}

// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.token))

// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return err
// 	}

// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusCreated {
// 		return fmt.Errorf("failed to push image: %s", res.Status)
// 	}
// 	return nil
// }
