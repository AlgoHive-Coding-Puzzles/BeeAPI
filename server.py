import sys
import signal
import atexit
from dotenv import load_dotenv
from app.auth import create_user_from_env
from app import app, loader

# Try to load environment variables from .env file if it exists
try:
    load_dotenv()
    create_user_from_env()
except ImportError:
    print("python-dotenv not installed, skipping .env file loading")

def on_exit():
    loader.unload()
    
def handle_signal(signum, frame):
    on_exit()
    sys.exit(0)

if __name__ == '__main__':
    loader.extract()
    loader.load()
    
    atexit.register(on_exit)
    
    signal.signal(signal.SIGTERM, handle_signal)
    signal.signal(signal.SIGINT, handle_signal)
    
    app.run(host='0.0.0.0')
    