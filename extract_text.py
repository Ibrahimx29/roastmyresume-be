import sys
import fitz  # PyMuPDF
import re
from typing import List, Tuple

def extract_text_from_pdf(pdf_file: str, detect_headers: bool = True) -> List[str]:
    """
    Extract text from a PDF file with improved formatting preservation using PyMuPDF.
    
    Args:
        pdf_file: Path to the PDF file
        detect_headers: Whether to attempt to detect and format headers
        
    Returns:
        List of strings containing text from each page
    """
    pdf_text = []
    
    try:
        # Open the PDF
        doc = fitz.open(pdf_file)
        
        for page_num in range(len(doc)):
            page = doc.load_page(page_num)
            
            # Extract text with layout preservation
            text = page.get_text("text")  # Simple text mode
            
            # If we want more layout control, we can use the blocks format
            if detect_headers:
                text = improve_header_formatting(page)
            
            # Clean up and normalize text
            text = clean_text(text)
            
            pdf_text.append(text)
        
        doc.close()
        return pdf_text
        
    except Exception as e:
        print(f"Error extracting text: {e}")
        return [f"Error: {str(e)}"]

def improve_header_formatting(page) -> str:
    """
    Extract text with improved header detection by analyzing text blocks.
    
    Args:
        page: A PyMuPDF page object
        
    Returns:
        Text with improved header formatting
    """
    # Get blocks which contain (x0, y0, x1, y1, "text", block_no, block_type)
    blocks = page.get_text("blocks")
    
    # Sort blocks by vertical position (top to bottom)
    blocks.sort(key=lambda b: b[1])  # Sort by y0 coordinate
    
    formatted_text = []
    last_y_bottom = -1
    line_spacing_threshold = 5  # Adjust based on your document
    
    for block in blocks:
        if block[6] == 0:  # Text blocks only (skip images)
            x0, y0, x1, y1 = block[:4]
            text = block[4]
            
            # Check if this might be a header based on:
            # 1. Distance from previous block
            # 2. Short length (headers tend to be shorter)
            # 3. Ends with a colon or contains title-like indicators
            
            is_likely_header = False
            
            # Check for significant gap from previous text
            if last_y_bottom > 0 and (y0 - last_y_bottom) > line_spacing_threshold:
                is_likely_header = True
                
            # Check for typical header patterns
            if re.match(r'^[A-Z][A-Za-z\s]{1,25}$', text.strip()):
                is_likely_header = True
                
            if re.search(r'(Summary|Education|Experience|Skills|Achievements|Certificates|Languages):', text):
                is_likely_header = True
            
            # Add extra line break before headers
            if is_likely_header and formatted_text:
                formatted_text.append("")
                
            formatted_text.append(text.strip())
            last_y_bottom = y1
    
    return "\n".join(formatted_text)
    
def clean_text(text: str) -> str:
    """
    Clean up and normalize extracted text.
    
    Args:
        text: Raw extracted text
        
    Returns:
        Cleaned text
    """
    # Replace problematic characters
    text = text.replace('\ufffd', '-')
    
    # Remove excessive whitespace while preserving paragraph breaks
    text = re.sub(r'\n\s*\n', '\n\n', text)
    text = re.sub(r' +', ' ', text)
    
    # Fix common layout issues
    # Fix merged lines (when name and contact info get merged)
    text = re.sub(r'([A-Za-z]+)(\+\d{2}\s+\d+)', r'\1\n\2', text)
    
    # Fix broken heading lines
    text = re.sub(r'([A-Za-z]+)(\d{1,2}/\d{4})', r'\1\n\2', text)
    
    return text

def format_extracted_text(pdf_text: List[str]) -> str:
    """
    Format the extracted text for better readability.
    
    Args:
        pdf_text: List of strings containing text from each page
        
    Returns:
        Formatted text as a single string
    """
    return "\n\n".join(pdf_text)

if __name__ == '__main__':
    if len(sys.argv) < 2:
        print("Usage: python pymupdf_extractor.py <PDF_FILE_PATH>")
        sys.exit(1)

    path = sys.argv[1]
    text_list = extract_text_from_pdf(path)
    formatted_text = format_extracted_text(text_list)
    print(formatted_text)