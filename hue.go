package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Hue struct {
	URL           string                 `json:"-"`
	Lights        map[string]Light       `json:"lights,omitempty"`
	Groups        map[string]Group       `json:"groups,omitempty"`
	Scenes        map[string]Scene       `json:"scenes,omitempty"`
	Sensors       map[string]interface{} `json:"sensors,omitempty"`
	Schedules     map[string]Schedule    `json:"schedules,omitempty"`
	Rules         map[string]interface{} `json:"rules,omitempty"`
	Resourcelinks map[string]interface{} `json:"resourcelinks,omitempty"`
	Config        map[string]interface{} `json:"config,omitempty"`
}

type Light struct {
	ID               string `json:"id,omitempty"`               // i.e. "1"
	Name             string `json:"name,omitempty"`             // i.e. "Light above chair"
	State            State  `json:"state,omitempty"`            //
	Type             string `json:"type,omitempty"`             // i.e. "Color temperature light"
	ModelID          string `json:"modelid,omitempty"`          // i.e. "LTW011"
	SWVersion        string `json:"swversion,omitempty"`        // i.e. "1.15.2_r19181"
	Manufacturername string `json:"manufacturername,omitempty"` // i.e. "Philips"
	Uniqueid         string `json:"uniqueid,omitempty"`         // i.e. "00:17:88:01:02:22:a1:da-0b"
	Productid        string `json:"productid,omitempty"`        // i.e. "Philips-LTW011-1-BR30CTv1"
}

type State struct {
	On             bool      `json:"on"`
	Bri            int       `json:"bri,omitempty"`            //LTW011, [1..254]
	ColorMode      string    `json:"colormode,omitempty"`      //i.e. "ct"
	CT             int       `json:"ct,omitempty"`             //LTW011. Measured in "reciprocal megakelvin": 500=2000K, 153=6500K
	Hue            int       `json:"hue,omitempty"`            //
	Effect         string    `json:"effect,omitempty"`         //
	Sat            int       `json:"sat,omitempty"`            //
	XY             []float32 `json:"xy,omitempty"`             //
	Alert          string    `json:"alert,omitempty"`          //i.e. "none"
	TransitionTime int       `json:"transitiontime,omitempty"` //
	Reachable      bool      `json:"reachable"`                // i.e. true
}

type Group struct {
	Name   string   `json:"name,omitempty"`
	Type   string   `json:"type,omitempty"`
	Class  string   `json:"class,omitempty"`
	Action Action   `json:"action,omitempty"`
	Lights []string `json:"lights,omitempty"`
}

type Action struct {
	CT        int    `json:"ct,omitempty"`
	Bri       int    `json:"bri,omitempty"`
	On        bool   `json:"on,omitempty"`
	Colormode string `json:"colormode,omitempty"`
	Alert     string `json:"alert,omitempty"`
}

type Scene struct {
	ID          string      `json:"id,omitempty"`
	Name        string      `json:"name,omitempty"`
	Lights      []string    `json:"lights"`
	Owner       string      `json:"owner"`
	Recycle     bool        `json:"recycle"`
	Locked      bool        `json:"locked"`
	LastUpdated string      `json:"lastupdated"`
	Version     int         `json:"version"`
	Appdata     interface{} `json:"appdata,omitempty"`
	Picture     string      `json:"picture,omitempty"`
}

type Schedule struct {
	Name        string  `json:"name"`
	Recycle     bool    `json:"recycle"`
	Description string  `json:"description"`
	Time        string  `json:"time"`
	Command     Command `json:"command"`
	Status      string  `json:"status"`
	Created     string  `json:"created"`
	Localtime   string  `json:"localtime"`
}

type Command struct {
	Body    interface{} `json:"body"`
	Address string      `json:"address"`
	Method  string      `json:"method"`
}

func Connect(url string) (*Hue, error) {
	contents, err := httpGet(url)
	if err != nil {
		return nil, err
	}
	var hue Hue
	err = json.Unmarshal(contents, &hue)
	if err != nil {
		return nil, err
	}
	hue.URL = url
	return &hue, nil
}

func main() {
	hueURL := os.Getenv("HUE_URL")
	if hueURL == "" {
		fmt.Println("Set the 'HUE_URL' environment variable to 'http://{host}/api/{user}'")
		os.Exit(1)
	}
	hue, err := Connect(hueURL)
	if err != nil {
		fmt.Println("Cannot read hub state:", err)
		os.Exit(1)
	}

	var lights []string
	for id, _ := range hue.Lights {
		lights = append(lights, id)
	}
	argv := os.Args
	argc := len(argv)
	if argc == 1 {
		usage()
	} else {
		switch argv[1] {
		case "light", "lights":
			if argc < 3 {
				for _, light := range hue.AllLights() {
					fmt.Printf("%s\t%d\t%dk\t%s\n", light.ID, light.Brightness(), light.Temperature(), light.Name)
				}
			} else {
				if argv[2] != "all" {
					lights = strings.Split(argv[2], ",")
					if len(lights) < 1 {
						usage()
					}
				}
				bright := "-"
				temp := "-"
				if argc > 3 {
					bright = argv[3]
					if argc > 4 {
						temp = argv[4]
					}
				}
				tempValue := ColorTempValue(temp)
				brightValue, on, ok := BrightnessValue(bright)
				if ok {
					hue.SetLights(lights, on, brightValue, tempValue)
				}
			}
		case "scene", "scenes":
			if argc < 3 {
				//list the valid scene names
				for id, scene := range hue.Scenes {
					fmt.Printf("%s\t%s\n", id, scene.Name)
				}
				os.Exit(0)
			}
			sceneName := strings.ToLower(argv[2])
			group := "1" //fix this later
			if argc > 3 {
				group = argv[3]
			}
			for sceneId, scene := range hue.Scenes {
				if strings.ToLower(scene.Name) == sceneName {
					hue.SetScene(sceneId, group)
					os.Exit(0)
				}
			}
			fmt.Printf("Cannot find scene %q\n", sceneName)
		default:
			usage()
		}
	}
}

func usage() {
	fmt.Println("usage: hue light")
	fmt.Println("usage: hue light <ids> on")
	fmt.Println("usage: hue light <ids> off")
	fmt.Println("usage: hue light <ids> bright")
	fmt.Println("usage: hue light <ids> bright temp")
	fmt.Println("usage: hue scene")
	fmt.Println("usage: hue scene <name>")
	fmt.Println(" bright - a value between 0 and 100, or 'on' or 'off' or a preset name like 'concentrate'. Default is '-' i.e. don't change")
	fmt.Println(" temp - a value between 2000 and 6500, or a symbol like 'warm' or 'cold', default is '-', i.e. don't change")
	fmt.Println(" light - an identifier for the light, i.e. '1' or '2', or a group name i.e. 'Office'. Default is 'all'")
	fmt.Println("with no arguments, the list of known lights is shown")
	os.Exit(1)
}

/*
func scene(name string) (string, string) {
	switch strings.ToLower(name) {
	case "relax":
		return "57", "2257"
	case "read":
		return "100", "2915"
	case "concentrate", "focus":
		return "100", "4347"
	case "energize":
		return "100", "6410"
	case "savanna sunset", "sunset", "savanna":
		return "79", "2617"
	case "tropical twilight", "tropical", "twilight":
		return "49", "3105"
	case "arctic aurora", "arctic", "aurora":
		return "54", "6535"
	case "spring blossom", "spring", "blossom":
		return "84", "4661"
	case "bright":
		return "100", "2732"
	case "dimmed":
		return "31", "2724"
	case "nightlight":
		return "1", "2257"
	default:
		return name, "-"
	}

}
*/

func (hue *Hue) SetLights(ids []string, on bool, bright int, temp int) {
	for _, id := range ids {
		err := hue.SetLight(id, on, bright, temp)
		if err != nil {
			fmt.Println("***", err)
			os.Exit(1)
		}
	}
}

func (hue *Hue) String() string {
	return pretty(hue)
}

func pretty(obj interface{}) string {
	b, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func httpGet(url string) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func httpPut(url string, body []byte) ([]byte, error) {
	//fmt.Printf("[%s]\n", url)
	bodyReader := strings.NewReader(string(body))
	request, err := http.NewRequest("PUT", url, bodyReader)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func (light Light) String() string {
	return pretty(light)
}
func (light Light) IsOn() bool {
	return light.State.On
}

func (light Light) Brightness() int {
	if !light.State.On {
		return 0
	}
	//bri is 1..254
	return (light.State.Bri * 100 / 255) + 1
}

func (light Light) Temperature() int {
	return 1000000 / light.State.CT
}

var tempSymbols = map[string]string{
	"glow":        "2000",
	"warm":        "2700",
	"white":       "3000",
	"bright":      "4000",
	"concentrate": "4347",
	"cool":        "5000",
	"cold":        "6500",
	"-":           "",
	"same":        "",
}

func ColorTempValue(temp string) int {
	if t, ok := tempSymbols[temp]; ok {
		if t == "" {
			return 0
		}
		temp = t
	} else if strings.HasSuffix(temp, "k") {
		temp = temp[:len(temp)-1]
	}
	tempInKiloKelvin, err := strconv.Atoi(temp)
	if err != nil {
		fmt.Printf("bad temp parameter: %q\n", temp)
		return 0 //don't change it
	}
	//expressed in kilokelvin i.e. "2700" for 2700k. Result is in mired (reciprocal megakelvin), which hue requires
	if tempInKiloKelvin < 2000 {
		return 500
	}
	if tempInKiloKelvin > 6500 {
		return 153
	}
	return 1000000 / tempInKiloKelvin
}

func BrightnessValue(bright string) (int, bool, bool) {
	on := true
	switch bright {
	case "off", "0":
		return 0, false, true
	case "on", "-", "same":
		return 0, true, true
	}
	brightnessFrom0to100, err := strconv.Atoi(bright)
	if err != nil {
		fmt.Printf("Bad bright argument: %q\n", bright)
		return 0, on, false
	}
	dim := brightnessFrom0to100 * 256 / 100
	if dim < 1 {
		return 1, on, true
	}
	if dim > 254 {
		return 254, on, true
	}
	return dim, on, true
}

func (hue *Hue) AllLights() []Light {
	lights := make([]Light, 0, len(hue.Lights))
	ids := make([]string, 0, len(lights))
	for k := range hue.Lights {
		ids = append(ids, k)
	}
	sort.Strings(ids)
	for _, id := range ids {
		light := hue.Lights[id]
		light.ID = id
		lights = append(lights, light)
	}
	return lights
}

func (hue *Hue) SetLight(id string, on bool, bright int, temp int) error {
	state := &State{On: on}
	if bright != 0 {
		state.Bri = bright
	}
	if temp != 0 {
		state.CT = temp
	}
	jsonState, err := json.Marshal(state)
	if err != nil {
		return err
	}
	url := hue.URL + "/lights/" + id + "/state"
	_, err = httpPut(url, jsonState)
	return err
}

func (scene Scene) String() string {
	return pretty(scene)
}

func (hue *Hue) AllScenes() ([]Scene, error) {
	url := hue.URL + "/scenes"
	contents, err := httpGet(url)
	if err != nil {
		return nil, err
	}
	scenesMap := map[string]Scene{}
	err = json.Unmarshal(contents, &scenesMap)
	if err != nil {
		return nil, err
	}
	scenes := make([]Scene, 0, len(scenesMap))
	for id, scene := range scenesMap {
		scene.ID = id
		scenes = append(scenes, scene)
	}
	return scenes, nil
}

func (hue *Hue) SetScene(sceneId string, group string) error {
	url := hue.URL + "/groups/" + group + "/action"
	state := map[string]string{"scene": sceneId}
	jsonState, err := json.Marshal(state)
	if err != nil {
		return err
	}
	_, err = httpPut(url, jsonState)
	return err
}
