resource "aws_lambda_layer_version" "lambda_layer" {
  layer_name          = "video_package_layer"
  s3_bucket           = aws_s3_bucket.lambda_bucket.id
  s3_key              = aws_s3_object.lambda_layer_zip.key

  compatible_runtimes = ["python3.11"]
}

resource "aws_lambda_function" "function" {
  s3_bucket = aws_s3_bucket.lambda_bucket.id
  s3_key    = aws_s3_object.lambda_bkt.key
  timeout =  900
  memory_size =  1020
  function_name    = var.lambda_name
  role             = aws_iam_role.lambda_role.arn
  handler          = "lambda_function.lambda_handler"
  runtime          = "python3.11"
  layers = [aws_lambda_layer_version.lambda_layer.arn]
  source_code_hash = data.archive_file.lambda_function.output_base64sha256
  ephemeral_storage {
    size = 1020
  }
}

data "archive_file" "lambda_function" {
  type        = "zip"
  source_dir  = "lambda/"
  output_path = "lambda.zip"
}
