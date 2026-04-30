@echo off
setlocal

cd /d "%~dp0"

where python >nul 2>nul
if errorlevel 1 (
  echo Python not found in PATH.
  exit /b 1
)

echo Starting core_agent_service from %cd%
python -m core_agent_service

endlocal
