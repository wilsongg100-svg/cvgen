curl -X POST http://localhost:8080/api/v1/resumes \
 -H "Content-Type: application/json" \
 -d @exemplo-cv.json \
 --output cv.pdf

curl -X POST http://localhost:8080/api/v1/resumes/import \
 -F "cv=@cv.pdf"
