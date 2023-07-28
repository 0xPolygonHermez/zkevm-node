
# Sequencer HTTP Endpoints Documentation

## Overview

The Sequencer uses an HTTP server with several dedicated endpoints to control its operational phases. This document outlines these endpoints, their functionality, and provides examples of request data and responses.

For security reasons, these endpoints are only accessible from the localhost. If remote access is necessary, ensure to establish a secure SSH tunnel. If the port needs to be exposed publicly for any reason, it is critical to protect it with stringent firewall rules, limiting access only to trusted sources.

## Endpoints

### 1. `/stopAfterCurrentBatch`

- **Method:** `POST`
- **Description:** Stops the Sequencer after it completes and closes the current batch.
- **Request Example:**
  ```
  POST /stopAfterCurrentBatch HTTP/1.1
  Host: localhost
  ```
- **Response Example:**
  ```json
  {
      "message": "Stopping after current batch"
  }
  ```

### 2. `/stopAtBatch`

- **Method:** `POST`
- **Description:** Stops the Sequencer after it completes the batch with the specified batch number.
- **Request Body:** JSON object containing `batchNumber` field.
- **Request Example:**
  ```
  POST /stopAtBatch HTTP/1.1
  Host: localhost
  Content-Type: application/json
  
  {
      "batchNumber": 5
  }
  ```
- **Response Example:**
  ```json
  {
      "message": "Stopping at specific batch"
  }
  ```

### 3. `/resumeProcessing`

- **Method:** `POST`
- **Description:** Resumes the Sequencer's operation from the batch following the last processed one before the stop command was issued.
- **Request Example:**
  ```
  POST /resumeProcessing HTTP/1.1
  Host: localhost
  ```
- **Response Example:**
  ```json
  {
      "message": "Resuming processing"
  }
  ```

### 4. `/getCurrentBatchNumber`

- **Method:** `GET`
- **Description:** Returns the number of the batch that is currently being processed by the Sequencer.
- **Request Example:**
  ```
  GET /getCurrentBatchNumber HTTP/1.1
  Host: localhost
  ```
- **Response Example:**
  ```json
  {
      "currentBatchNumber": "3"
  }
  ```
