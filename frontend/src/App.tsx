// Dependencies
import { useState } from "react";
import logo from "./logo.svg";
import "./App.css";
import {
  useQuery,
  useMutation,
  useQueryClient,
  QueryClient,
  QueryClientProvider,
} from "@tanstack/react-query";
import axios from "axios";

// Create a client
const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <CardUpgrades />
    </QueryClientProvider>
  );
}

function CardUpgrades(): JSX.Element {
  const { isLoading, error, data, isFetching } = useQuery({
    queryKey: ["cardUpgrades"],
    queryFn: () =>
      axios({
        method: "get",
        url: "http://localhost:8080/api/cards/upgrades",
        headers: {
          "Access-Control-Allow-Origin": "*",
        },
      }).then((res) => {
        return res.data;
      }),
  });

  if (isLoading) return <div>"Loading..."</div>;

  if (error) return <div>"An error has occurred: " + error.message</div>;

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          {isFetching
            ? "Fetching your cards..."
            : data.map((card: string) => {
                return <p>{card}</p>;
              })}
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
      </header>
    </div>
  );
}

export default App;
