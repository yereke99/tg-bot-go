# Use the official Python image as a base image
FROM python:3.8.0

# Set the working directory in the container
WORKDIR /app

# Copy the local code to the container image
COPY . /app

# Install dependencies
RUN pip install --no-cache-dir \
    aiogram==2.24 \
    aiohttp==3.8.6 \
    aioschedule==0.5.2 \
    fpdf==1.7.2 \
    frozenlist==1.4.1 \
    idna==3.10 \
    magic-filter==1.0.12 \
    multidict==6.1.0 \
    pillow==10.4.0 \
    pypng==0.20220715.0 \
    qrcode==7.4.2 \
    reportlab==4.2.2 \
    setuptools==75.1.0 \
    typing_extensions==4.12.2 \
    wheel==0.44.0 \
    yarl==1.11.1

# Expose the port the app runs on
EXPOSE 8080

# Run the application
CMD ["python", "run.py"]