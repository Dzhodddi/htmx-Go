from pydantic import BaseModel, Field
from langchain_core.tools import BaseTool
import os
import re

class CodebaseGrepInput(BaseModel):
    search_directory: str = Field(description="Directory to search")
    pattern: str = Field(description="Regex pattern to search for")
    file_extension: str | None = Field(description="Optional file extension filter", default=None)
    case_sensitive: bool = Field(description="Case-sensitive search flag", default=False)

class CodebaseGrepTool(BaseTool):
    name: str = "codebase_grep"
    description: str = "Searches codebase for regex pattern with optional file extension filter"
    args_schema: type[BaseModel] = CodebaseGrepInput

    def _run(self, search_directory: str, pattern: str, file_extension: str | None = None, case_sensitive: bool = False) -> str:
        try:
            results = []
            flags = 0 if case_sensitive else re.IGNORECASE
            excluded_dirs = {'.git', 'node_modules', 'venv', '.venv', '__pycache__'}

            for root, dirs, files in os.walk(search_directory):
                dirs[:] = [d for d in dirs if d not in excluded_dirs]

                for file in files:
                    if file_extension and not file.endswith(file_extension):
                        continue

                    file_path = os.path.join(root, file)
                    try:
                        with open(file_path, 'r', encoding='utf-8') as f:
                            content = f.read()
                            matches = re.findall(pattern, content, flags=flags)
                            if matches:
                                results.append(f"{file_path}:")
                                for match in matches[:5]:  # Limit 5 matches per file
                                    results.append(f"  - {match}")
                                if len(matches) > 5:
                                    results.append("  ... (truncated)")
                    except (UnicodeDecodeError, PermissionError):
                        pass  # Skip binary files and permission issues

            if len(results) > 100:
                results = results[:100] + ["... (100 match limit reached)"]

            return '\n'.join(results) or 'No matches found'
        except Exception as e:
            return f"Error executing codebase_grep: {str(e)}"
