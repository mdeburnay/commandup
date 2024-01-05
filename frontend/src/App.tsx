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
  console.log("App");
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
      }).then((res) => {
        console.log(res.data);
        return res.data;
      }),
  });

  if (isLoading) return <div>"Loading..."</div>;

  if (error) {
    return <div>{error.toString()}</div>;
  }

  if (data) console.log(data);

  return (
    <div className="App">
      <header className="App-header">
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          {isFetching
            ? "Fetching your cards..."
            : data.matchingCards.map((card: string, index: number) => {
                return <div key={index}>{card}</div>;
              })}
        </a>
      </header>
    </div>
  );
}

export default App;
