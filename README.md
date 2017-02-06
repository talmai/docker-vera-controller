### Docker container for RaspberryPi which allow for Veralite™ Smart Home Controller Automation

Use Amazon Echo to voice control your home automation devices through http commands sent to your home automation controller or built-in direct control of the Harmony Hub and Nest. 

Also includes MQTT server.

This is a docker container for bwssystems' [ha-bridge](https://github.com/bwssytems/ha-bridge), forked from aptalca's [docker-ha-bridge](https://github.com/aptalca/docker-ha-bridge) and merged with work from hypriot's [rpi-java](https://github.com/hypriot/rpi-java).

#### Installation instruction:

You can run this docker with the following command:

```docker run -d --name="Home-Automation-Bridge" --net="host" -e SERVERIP="192.168.X.X" -e SERVERPORT="XXXX" -v /path/to/config/:/config:rw -v /etc/localtime:/etc/localtime:ro aptalca/home-automation-bridge```

- Replace the SERVERIP variable (192.168.X.X) with your server's IP
- Replace the SERVERPORT variable (XXXX) with whichever port you choose for the web gui.
- Replace the "/path/to/config" with your choice of location
- If the `-v /etc/localtime:/etc/localtime:ro` mapping causes issues, you can try `-e TZ="<timezone>"` with the timezone in the format of "America/New_York" instead

##### Optional Variables for the run command
- By default, this will install the latest version on bwssystems github repo, and will auto update itself to the latest version on each container start, but if you want to run a different version (to go back to the previous version perhaps), include the following environment variable in your docker run command `-e VERSION="X.X.X"`
- Once installed, open the WebUI at `http://SERVERIP:SERVERPORT/` and enter your Vera, Harmony and Nest info.
  
### Other code in this repo

`hue.go`: First, get a client "user" value, then set the `HUE_URL` environment variable to `http://{host}/api/{user}`.

	$ hue
	usage: hue light
	usage: hue light <ids> on
	usage: hue light <ids> off
	usage: hue light <ids> bright
	usage: hue light <ids> bright temp
	usage: hue scene
	usage: hue scene <name>
	 bright - a value between 0 and 100, or 'on' or 'off' or a preset name like 'concentrate'. Default is '-' i.e. don't change
	 temp - a value between 2000 and 6500, or a symbol like 'warm' or 'cold', default is '-', i.e. don't change
	 light - an identifier for the light, i.e. '1' or '2', or a group name i.e. 'Office'. Default is 'all'
	with no arguments, the list of known lights is shown

	$ hue lights
	1	0	2732k	Light above chair
	2	0	2732k	Light above desk
	3	100	3508k	Desk lamp
	4	100	3508k	Reading lamp
	$ hue lights 1,2 on
	$ hue lights
	1	100	2732k	Light above chair
	2	100	2732k	Light above desk
	3	100	3508k	Desk lamp
	4	100	3508k	Reading lamp
	$ hue lights 1,2 on 3000k
	$ hue lights
	1	100	3003k	Light above chair
	2	100	3003k	Light above desk
	3	100	3508k	Desk lamp
	4	100	3508k	Reading lamp
	$ hue lights 1,2 50 2800k
	1	51	2801k	Light above chair
	2	51	2801k	Light above desk
	3	100	3508k	Desk lamp
	4	100	3508k	Reading lamp
	$ hue scene
	HkUBtacEYgNmiAt	Energize
	7ZczJxXK8o6882B	Arctic aurora
	EZD73K3Q9DsvIEj	Bright
	9DCUzOQJpWKxE2b	Tropical twilight
	CL-v3NHEez2y2aw	Concentrate
	OQIbZPuQ0cRwaIH	Morning work
	6XpxObjpFbT-zhO	Relax
	azoXeiq8xynwfYS	Read
	j3pErGeKZL40zDL	Nightlight
	56L9MVO5u2OWuJb	Savanna sunset
	TPBNmbLsrHCP3Iu	Spring blossom
	ttN1Dw3J7xbCcsu	Video
	pec4GCvLs43jCA5	Energize
	WTLiYS5-b4hOM0C	Dimmed
	$ hue scene 'tropical twilight'
	$ hue lights
	1	49	3105k	Light above chair
	2	49	3105k	Light above desk
	3	49	3105k	Desk lamp
	4	49	3105k	Reading lamp
