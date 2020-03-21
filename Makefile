.PHONY: deploy delete

deploy:
	gcloud functions deploy subscriptions$(commit) --entry-point Subscribe --runtime go113 --trigger-http --memory 128MB --region $(region) --service-account $(serviceaccount)

delete:
	gcloud functions delete subscriptions$(commit) --region $(region)

grant_permission:
	gcloud functions add-iam-policy-binding subscription$(commit) --region $(region) \
   --member "serviceAccount:$(serviceaccount)" \
   --role "roles/cloudfunctions.invoker" \
   --project $(project)

test:
	go test .

integration: # TODO add integration tests
	go test .

smoke: # TODO add smoke tests
	go test .