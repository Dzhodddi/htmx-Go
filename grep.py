import os
import re

def CodebaseGrepTool(search_directory, pattern, file_extension=None, case_sensitive=False):
    """Search codebase for regex pattern with optional file extension filter"""
    results = []
    ignore_dirs = {'.git', 'node_modules', 'venv', '.venv', '__pycache__'}
    flags = 0 if case_sensitive else re.IGNORECASE
n    
    for root, dirs, files in os.walk(search_directory):
        # Remove ignored directories
        dirs[:] = [d for d in dirs if os.path.join(root, d) not in ignore_dirs]
        
        for file in files:
            if file_extension and not file.endswith(file_extension):
                continue
            
            file_path = os.path.join(root, file)
            try:
                with open(file_path, 'r', encoding='utf-8') as f:
                    content = f.read()
                    matches = re.finditer(pattern, content, flags)
                    
                    for match in matches:
                        results.append({
                            'file': file_path,
                            'match': match.group(),
                            'line': content.splitlines()[match.start().count('\n')]
                        })
                        
                # Stop if we've reached 100 matches
                if len(results) >= 100:
                    break
            except (UnicodeDecodeError, PermissionError):
                # Skip binary files and access-denied files
                continue
    
    return results[:100]