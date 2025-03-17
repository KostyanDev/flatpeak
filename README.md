# **Carbon Intensity Optimization API**

## **Project Description**
This microservice is designed to analyze and provide optimal time slots when carbon intensity is at its lowest. Carbon intensity represents the amount of CO₂ emitted per kWh, and the system helps identify the most environmentally friendly time slots for energy consumption.
Built with Golang, this RESTful API fetches carbon intensity forecast data, processes it efficiently using optimized algorithms, and provides users with the best time slots based on their requested duration.

## **Project Goals and Approach**
The objective of this project was to build a high-performance, scalable API that:
-	Efficiently processes carbon intensity data.
-	Handles concurrency and optimized calculations to determine the best slots.
-	Supports both continuous and non-continuous time slots for flexible use cases.
-	Ensures robustness with proper input validation, error handling, and logging.
-	Follows clean architecture principles, separating concerns between fetching data, business logic, and API handling.
---

## **Requirements**
-   Fetch Carbon Intensity Forecast
-   Retrieve carbon intensity forecast for the next 24 hours.
-   Process the forecast data to determine optimal time slots.
-   Identify the lowest carbon intensity period based on user-defined duration.
-   Support both continuous (consecutive periods) and non-continuous (dispersed optimal periods) selection.
-   If a requested slot spans multiple 30-minute intervals, return a weighted average of intensity values.
-   Provide an API endpoint to retrieve optimal time slots.
-   Accept and validate query parameters (duration, continuous).
-   Return data in JSON format.
---
## **Makefile Commands**
| Command          | Description                                      |
|-----------------|--------------------------------------------------|
| `make all`      | Build the application and start all containers   |
| `make build`    | Build the application binary                     |
| `make up`       | Start the application and database               |
| `make down`     | Stop the application and database                |
| `make restart`  | Restart the application                          |
| `make test`     | Run all integration tests                        |
| `make swagger`  | Generate Swagger documentation                   |

---
## **Usage** 

After starting the application, the API will be available at:
http://0.0.0.0:8080
Swagger documentation can be accessed at:
http://0.0.0.0:8080/swagger/index.html
### **API Endpoints**
#### **URL**
| Method | Endpoint | Description                             |
|--------|----------|-----------------------------------------|
| `GET`  | `/slots` | Retrieves the best time slots with the lowest carbon intensity for the given duration|

---