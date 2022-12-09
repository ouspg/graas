FROM python:3-slim
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd -g 1000 appuser && \
    useradd -u 1000 -g appuser -m -s /sbin/nologin appuser

WORKDIR /home/appuser/app
COPY secret/secret /usr/share/
USER appuser
COPY bob/requirements.txt /home/appuser/app/requirements.txt
RUN pip install -r requirements.txt
# Copy our code from the current folder to /app inside the container
COPY bob/message.py .

CMD ["python", "message.py"]