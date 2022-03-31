build_sensor_image:
	docker build ./cmd/sensor -t youla_dev_internship_task_go_sensor:latest

run_sensors: build_sensor_image
	docker-compose up -d
