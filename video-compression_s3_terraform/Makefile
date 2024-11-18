.PHONY: build clean

LAYER_ZIP=python.zip
build: lambda/lambda_function.py
	zip -r9 lambda.zip lambda/lambda_function.py

layer: lambda/
	pip3 install -r lambda/requirements.txt -t python
	zip -r $(LAYER_ZIP) python
clean:
	rm -f lambda.zip && rm -rf python && rm python.zip


init:
	terraform init

plan:
	terraform plan -var-file=vars.tfvars


apply:
	terraform apply -var-file=vars.tfvars --auto-approve

dis:
	terraform destroy -var-file=vars.tfvars --auto-approve