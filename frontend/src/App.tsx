// Dependencies
import { useState, useEffect } from "react";
import "./App.css";
import {
  useQuery,
  useMutation,
  QueryClient,
  QueryClientProvider,
} from "@tanstack/react-query";
import axios from "axios";
import Papa from "papaparse";

// Create a client
const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <CardUpgrades />
      <CardUploadButton />
    </QueryClientProvider>
  );
}

function CardUploadButton(): JSX.Element {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files) {
      setSelectedFile(event.target.files[0]);
    }
  };

  const handleFileUpload = () => {
    if (selectedFile) {
      const formData = new FormData();
      formData.append("file", selectedFile);

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
          console.log(response);
          // Handle the response from your backend here
          console.log("File uploaded successfully:", response.data);
        })
        .catch((error) => {
          // Handle any errors that occur during the upload
          console.error("Error uploading file:", error);
        });
    }
  };

  return (
    <>
      <input type="file" onChange={handleFileChange} />
      <button style={{ height: 50, width: 200 }} onClick={handleFileUpload}>
        Upload File
      </button>
    </>
  );
}

function CardUpgrades(): JSX.Element {
  const { error, data } = useMutation({
    mutationFn: () =>
      axios({
        method: "POST",
        url: "http://localhost:8080/api/cards/upgrades",
      }).then(({ data }) => {
        return data;
      }),
  });

  if (error) {
    return <div>{error.toString()}</div>;
  }

  return (
    <div className="App">
      <>
        {data &&
          data.map(
            (
              { title, cards }: { title: string; cards: string[] },
              i: number
            ) => {
              return (
                <div key={i}>
                  <h2 style={{ fontSize: 16 }}>{title}</h2>
                  <div>
                    {cards.map((card: string, i: number) => {
                      return (
                        <div key={i} style={{ fontSize: 14 }}>
                          {card}
                        </div>
                      );
                    })}
                  </div>
                </div>
              );
            }
          )}
      </>
    </div>
  );
}

export default App;
