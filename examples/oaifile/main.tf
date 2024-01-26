terraform {
  required_providers {
    oaifile = {
      source = "app.terraform.io/onelittlenightmusic/oaifile"
      version = "0.0.1"
    }
    # aws = {
    #   source  = "hashicorp/aws"
    #   version = "~> 3.0"  # Replace with your desired version
    # }
  }
}

provider "oaifile" {
  host     = "http://localhost:5001"
}

# provider "aws" {
#   region = "us-west-1" # Change to your AWS region
# }

# data "aws_s3_bucket_object" "example1" {
#   bucket = "blog-documents"
#   key    = "graphQLschemafunc.pdf"
# }

resource "oaifile_file" "example" {
  filepath = "s3://blog-documents/graphQLschemafunc.pdf"
  name = "test"
}

# output "oaifile_file" {
#   value = aws_s3_bucket_object.example1
# }
