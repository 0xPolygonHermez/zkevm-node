# main.tf

provider "aws" {
  region = "us-east-1"  # Change this to your desired AWS region
}

# Create a VPC
resource "aws_vpc" "my_vpc" {
  cidr_block = "10.0.0.0/16"
}

# Create a public subnet
resource "aws_subnet" "public_subnet" {
  vpc_id                  = aws_vpc.my_vpc.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = "us-east-1a"  # Change this to your desired AZ
  map_public_ip_on_launch = true
}

# Create a security group allowing inbound traffic on port 22 and 8545 (Geth RPC)
resource "aws_security_group" "allow_ssh_and_rpc" {
  name        = "allow_ssh_and_rpc"
  description = "Allow SSH and RPC traffic"

  vpc_id = aws_vpc.my_vpc.id

  ingress {
    from_port = 22
    to_port   = 22
    protocol  = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port = 8545
    to_port   = 8545
    protocol  = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Create an EC2 instance
resource "aws_instance" "geth_node" {
  ami           = "ami-xxxxxxxxxxxxxxxxx" 
  instance_type = "t2.micro"
  key_name      = "your-key-pair-name" 

  subnet_id        = aws_subnet.public_subnet.id

  user_data = <<-EOF
              #!/bin/bash
              sudo apt-get update
              sudo apt-get install -y software-properties-common
              sudo add-apt-repository -y ppa:ethereum/ethereum
              sudo apt-get update
              sudo apt-get install -y geth
              EOF

  tags = {
    Name = "geth-node"
  }
}
