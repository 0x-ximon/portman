provider "aws" {
  region = "us-west-2"
}

variable "env" {
  description = "Environment name"
  type        = string
}

variable "name" {
  description = "Name of the resource"
  type        = string
}

resource "aws_instance" "example" {
  ami           = "ami-0c94855ba95c71c99"
  instance_type = "t2.micro"
  tags = {
    Name = "${var.env}-${var.name}"
  }
}