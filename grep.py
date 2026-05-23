from pydantic import BaseModel, Field
from langchain_core.tools import BaseTool
import os
import re

class CodebaseGrepToolInput(BaseModel):
    search_directory: str = Field(description="Directory to search for the pattern")
    pattern: str = Field(description="Regex pattern to search for")
    file_extension: str = Field(description="File extension to filter by", default=None)
    case_sensitive: bool = Field(description="Whether the search should be case sensitive", default=False)

class CodebaseGrepTool(BaseTool):
    name: str = "codebase_grep_tool"
    description: str = "Search a given directory for a regex pattern"
    args_schema: type[BaseModel] = CodebaseGrepToolInput
    
    def _run(self, search_directory: str, pattern: str, file_extension: str = None, case_sensitive: bool = False) -> str:
        try:
            # Initialize results list
            results = []
            
            # Compile the regex pattern
            if case_sensitive:
                regex = re.compile(pattern)
            else:
                regex = re.compile(pattern, re.IGNORECASE)
            
            # Walk through the directory
            for root, dirs, files in os.walk(search_directory):
                # Ignore certain directories
                if '.git' in dirs:
                    dirs.remove('.git')
                if 'node_modules' in dirs:
                    dirs.remove('node_modules')
                if 'venv' in dirs:
                    dirs.remove('venv')
                if '.venv' in dirs:
                    dirs.remove('.venv')
                if '__pycache__' in dirs:
                    dirs.remove('__pycache__')
                
                # Iterate over files
                for file in files:
                    # Check file extension if provided
                    if file_extension and not file.endswith(file_extension):
                        continue
                    
                    # Construct file path
                    file_path = os.path.join(root, file)
                    
                    # Open and read the file
                    try:
                        with open(file_path, 'r') as f:
                            content = f.read()
                    except UnicodeDecodeError:
                        # Silently skip binary files
                        continue
                    
                    # Search for the pattern
                    matches = regex.findall(content)
                    
                    # Add matches to results
                    for match in matches:
                        results.append(f"{file_path}: {match}")
                        
                    # Cap results at 100 matches
                    if len(results) >= 100:
                        break
                
                # Cap results at 100 matches
                if len(results) >= 100:
                    break
            
            # Return results as a string
            return '\n'.join(results)
        except Exception as e:
            return f"Error executing codebase_grep_tool: {str(e)}"
