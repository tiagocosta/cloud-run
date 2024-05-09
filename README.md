# cloud-run
This is a simple service running o Google Cloud Run. It receives a zip code and returns the temperature of a location in Celsius, Fahrenheit and Kelvin.

The main idea of this project was, actually, test a simple and effective way of deploying an app on Google Cloud Run.

## API URL
https://cloud-run-goexpert-yt254bdq3q-uc.a.run.app/weather

## Usage
curl -XPOST -H 'Content-Type: application/json' -d '{"zipcode":"71218010"}' 'https://cloud-run-goexpert-yt254bdq3q-uc.a.run.app/weather'

## Dependencies
Go 1.22.2
