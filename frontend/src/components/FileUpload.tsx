// Dependencies
import { useRef } from "react";
import axios from "axios";
import { Button } from "./Button";

export function FileUpload(): JSX.Element {
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleButtonClick = () => {
    fileInputRef.current?.click();
  };

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
    <>
      <Button text="Upload Collection" onClick={handleButtonClick} />
      <input
        type="file"
        className="hidden"
        onChange={handleFileChange}
        ref={fileInputRef}
      />
    </>
  );
}
