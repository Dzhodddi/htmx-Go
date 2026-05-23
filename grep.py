import os
import re
from pydantic import BaseModel, Field
from langchain_core.tools import BaseTool
class CodebaseGrepToolInput(BaseModel):
    search_directory: str = Field(description="The directory to search for the pattern.")
    pattern: str = Field(description="The regex pattern to search for.")
    file_extension: Optional[str] = Field(description="The file extension to search for (optional).")
    case_sensitive: bool = Field(default=False, description="Whether the search should be case sensitive.")
class CodebaseGrepTool(BaseTool):
    name: str = "codebase_grep_tool"
    description: str = "Searches a given directory for a regex pattern using Python's re and os modules."
    args_schema: type[BaseModel] = CodebaseGrepToolInput

    def _run(self, search_directory: str, pattern: str, file_extension: Optional[str] = None, case_sensitive: bool = False) -> str:
        try:
            ignored_directories = ['.git', 'node_modules', 'venv', '.venv', '__pycache__']
            results = []
            for root, dirs, files in os.walk(search_directory):
                for file in files:
                    if file_extension and not file.endswith(file_extension):
                        continue
                    file_path = os.path.join(root, file)
                    try:
                        with open(file_path, 'r', encoding='utf-8') as f:
                            content = f.read()
                            if re.search(pattern, content, re.MULTILINE | (not case_sensitive and re.IGNORECASE)):
                                results.append(file_path)
                    except UnicodeDecodeError:
                        pass
                    if len(results) >= 100:
                        break
            return str(results)
        except Exception as e:
            return f"Error executing codebase_grep_tool: {str(e)}"
