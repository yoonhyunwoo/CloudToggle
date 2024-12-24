# How to install
Follow these steps to set up and run CloudToggle on your local or production environment.

---

## **🛠️ Prerequisites**

Before you begin, ensure you have the following tools installed:

1. **Go** (version 1.18 or higher): [Download Go](https://golang.org/dl/)
2. **Docker**: [Download Docker](https://www.docker.com/get-started)
3. **PostgreSQL** (if not using Docker): [Install PostgreSQL](https://www.postgresql.org/download/)
4. **Git**: [Download Git](https://git-scm.com/)

---

## **📦 Step 1: Clone the Repository**

Clone the CloudToggle repository from GitHub to your local machine:

```bash
git clone https://github.com/your-organization/cloudtoggle.git
cd cloudtoggle
```

---

## **🐳 Step 2: Set Up the Database with Docker**

CloudToggle requires a PostgreSQL database for storing resource group and schedule information. Set it up using Docker:

```bash
make up
```

This command will:
- Launch a PostgreSQL container.
- Apply database migrations to create the necessary schema.

To stop and clean up the database container:

```bash
make down
```

---

## **⚙️ Step 3: Configure Environment Variables**

Create a `.env` file in the project root directory and configure the following variables:

```plaintext
DB_URL=postgres://postgres:postgres@localhost:5432/cloudtoggle?sslmode=disable
JWT_SECRET=your_jwt_secret_key
AWS_ACCESS_KEY_ID=your_aws_access_key
AWS_SECRET_ACCESS_KEY=your_aws_secret_key
AWS_REGION=your_aws_region
```

Replace `your_...` placeholders with your actual credentials.

---

## **🚀 Step 4: Run the Application**

Start the CloudToggle server locally:

```bash
make run
```

The application will be available at **http://localhost:8080**.

---

## **🔗 Step 5: Test the API**

Use a tool like **Postman** or **cURL** to test the API. Start by logging in to receive a JWT token:

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "your_password"}'
```

---

## **🧪 Step 6: Build the Application**

To create a production-ready binary:

```bash
make build
```

This will generate a `cloudtoggle` binary in the project root directory.

---

## **🎉 Congratulations**

You’ve successfully set up CloudToggle! For further details, refer to:

- [API Documentation](../api)
- [Scheduling Guide](scheduling.md)