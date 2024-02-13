// Dependencies
import axios from "axios";

export function FileUpload(): JSX.Element {
  const handleFileUpload = (file: File) => {
    const formData = new FormData();
    formData.append("file", file);

    axios
      .post(
        "http://localhost:8080/api/cards/upload-card-collection",
        formData,
        {
          headers: {
            "Content-Type": "multipart/form-data",
          },
        }
      )
      .then((response) => {
        console.log("File uploaded successfully:", response.data);
      })
      .catch((error) => {
        console.error("Error uploading file:", error);
      });
  };

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files) {
      const file = event.target.files[0];
      handleFileUpload(file); // Upload the file as soon as it's selected
    }
  };

  return (
    <label className="w-64 flex flex-col items-center px-2 py-2 bg-white text-blue rounded-lg shadow-lg tracking-wide uppercase border border-blue cursor-pointer hover:bg-blue hover:text-white">
      <span className="mt-2 text-xs text-center">
        Upload Your Card Collection
      </span>
      <input type="file" className="hidden" onChange={handleFileChange} />
    </label>
  );
}
