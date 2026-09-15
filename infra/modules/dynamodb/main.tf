resource "aws_dynamodb_table" "analytics" {
  name         = var.table_name
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "flag_id"
  range_key    = "event_timestamp"

  attribute {
    name = "flag_id"
    type = "S"
  }

  attribute {
    name = "event_timestamp"
    type = "S"
  }

  point_in_time_recovery {
    enabled = true
  }

  server_side_encryption {
    enabled = true
  }

  tags = var.tags
}
