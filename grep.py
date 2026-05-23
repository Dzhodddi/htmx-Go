import os
import re

def codebase_grep(search_directory: str, pattern: str, file_extension: str = None, case_sensitive: bool = False) -> list:
    matches = []
    flags = 0 if case_sensitive else re.IGNORECASE
    
    for root, dirs, files in os.walk(search_directory):
        # Skip ignored directories
        if any(ignore_dir in dirs for ignore_dir in ['.git', 'node_modules', 'venv', '.venv', '__pycache__']):
            dirs[:] = [d for d in dirs if d not in ['.git', 'node_modules', 'venv', '.venv', '__pycache__']]
            continue
        
        for file in files:
            if file_extension and not file.endswith(file_extension):
                continue
            
            file_path = os.path.join(root, file)
            try:
                with open(file_path, 'r', encoding='utf-8') as f:
                    content = f.read()
                    if re.search(pattern, content, flags):
                        matches.append(file_path)
                        if len(matches) >= 100:
                            return matches
            except UnicodeDecodeError:
                # Skip binary files
                continue
    return matches