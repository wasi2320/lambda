resource "aws_dynamodb_table" "video_processing_table" {
  name           =  var.db_table_name
  billing_mode   = "PROVISIONED"
  read_capacity  = 1
  write_capacity = 1
  hash_key =  "object_key"
  range_key = "processed"
  attribute {
    name = "object_key"
    type = "S"
  }

  attribute {
    name = "processed"
    type = "N"
  }

  attribute {
    name = "error"
    type = "S"
  }
    global_secondary_index {
    name               = "processedIndex"
    hash_key           = "processed"
    range_key          = "error"
    write_capacity     = 1
    read_capacity      = 1
    projection_type    = "INCLUDE"
    non_key_attributes = ["object_key"]
  }
}
