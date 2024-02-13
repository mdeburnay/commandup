// Packages
import { useRef } from "react";
import { useMutation } from "@tanstack/react-query";
import axios from "axios";

// Components
import { Button } from "./Button";

export function FileUpload(): JSX.Element {
  const fileInputRef = useRef<HTMLInputElement>(null);

  const mutation = useMutation({
    mutationFn: async (file: File) => {
      const formData = new FormData();
      formData.append("file", file);

      return axios.post(
        "http://localhost:8080/api/cards/upload-card-collection",
        formData,
        {
          headers: {
            "Content-Type": "multipart/form-data",
          },
        }
      );
    },
  });

  const handleButtonClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files) {
      const file = event.target.files[0];
      mutation.mutate(file);
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
      {mutation.isPending && <div>Uploading...</div>}
      {mutation.isError && (
        <div>
          An error occurred:{" "}
          {mutation.error instanceof Error
            ? mutation.error.message
            : "Unknown error"}
        </div>
      )}
      {mutation.isSuccess && <div>File uploaded successfully!</div>}
    </>
  );
}
