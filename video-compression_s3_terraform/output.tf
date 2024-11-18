output "lambda_function_name" {
  value = aws_lambda_function.function.function_name
}

output "invoke_url" {
  value = "${aws_api_gateway_deployment.deployment.invoke_url}/${aws_api_gateway_resource.resource.path_part}"
}