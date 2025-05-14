import sys
import PyPDF2

def extract_text_from_pdf(pdf_file: str) -> [str]: # type: ignore
    with open(pdf_file, 'rb') as pdf:
        reader = PyPDF2.PdfReader(pdf, strict=False)
        pdf_text = []

        for page in reader.pages:
            content = page.extract_text()
            pdf_text.append(content)

        return pdf_text

if __name__ == '__main__':
    if len(sys.argv) < 2:
        print("Usage: python extract_text.py <PDF_FILE_PATH>")
        sys.exit(1)

    path = sys.argv[1]
    text = extract_text_from_pdf(path)
    print("\n---PAGE---\n".join(text))
