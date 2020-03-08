.PHONY: deploy delete

deploy:
	gcloud functions deploy subscribe$(commit) --entry-point Subscribe --runtime go113 --trigger-http --memory 128MB --region $(region) --serviceaccount $(serviceaccount)

delete:
	gcloud functions delete subscribe$(commit) --region europe-west1

test:
	go test .

integration: # TODO add integration tests
	go test .

smoke: # TODO add smoke tests
	go test .