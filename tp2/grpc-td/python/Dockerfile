FROM python:3.14-slim
WORKDIR /app
RUN pip install --no-cache-dir "grpcio>=1.85.1" "protobuf>=7.35.1"
COPY *.py ./
ENV PYTHONUNBUFFERED=1
EXPOSE 50051
CMD ["python", "server.py"]
