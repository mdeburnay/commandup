// Dependencies
import {
  useMutation,
  QueryClient,
  QueryClientProvider,
} from "@tanstack/react-query";
import axios from "axios";

// Styles
import "./index.css";

// Components
import { FileUpload } from "./components/FileUpload";

// Create a client
const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <CardUpgrades />
      <FileUpload />
    </QueryClientProvider>
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
