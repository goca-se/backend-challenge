from fastapi import Depends, HTTPException, status
from fastapi.security import HTTPBasic, HTTPBasicCredentials
import secrets
import base64

security = HTTPBasic()

# Hardcoded credentials for the mock API
# In a real scenario, these would be stored securely, possibly as environment variables
VALID_USERNAME = "entregas"
VALID_PASSWORD = "rapidas"


def verify_credentials(credentials: HTTPBasicCredentials = Depends(security)):
    """
    Verify the HTTP Basic Auth credentials.
    """
    is_username_valid = secrets.compare_digest(credentials.username, VALID_USERNAME)
    is_password_valid = secrets.compare_digest(credentials.password, VALID_PASSWORD)
    
    if not (is_username_valid and is_password_valid):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid credentials",
            headers={"WWW-Authenticate": "Basic"},
        )
    
    return credentials.username


# For reference, the encoded auth header value would be:
# "Basic ZW50cmVnYXM6cmFwaWRhcw=="
# which is base64 encoded "entregas:rapidas" 