FROM python:3-slim
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    curl \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd -g 1000 appuser && \
    useradd -u 1000 -g appuser -m -s /sbin/nologin appuser

WORKDIR /home/appuser/app
USER appuser
COPY requirements.txt /home/appuser/app/requirements.txt
RUN pip install -r requirements.txt
# Copy our code from the current folder to /app inside the container
COPY app.py .
ENV PATH=$PATH:/home/appuser/.local/bin/
EXPOSE 80

CMD ["gunicorn", "app:app", "-b", "0.0.0.0:80", "--log-file", "-", "--access-logfile", "-", "--workers", "4", "--keep-alive", "0"]