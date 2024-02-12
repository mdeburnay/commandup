// Dependencies
import { useQuery } from "@tanstack/react-query";
import axios from "axios";

// Components
import { FileUpload } from "../components/FileUpload";

export const Home = () => {
  return (
    <>
      <FileUpload />
      <CardUpgrades />
    </>
  );
};

function CardUpgrades(): JSX.Element {
  const { error, data } = useQuery({
    queryKey: ["card-upgrades"],
    queryFn: () =>
      axios({
        method: "GET",
        url: "http://localhost:8080/api/cards/upgrades",
      }).then(({ data }) => {
        return data;
      }),
  });

  if (error) {
    return <div>{error.toString()}</div>;
  }

  return (
    <div className="w-full justify-evenly flex flex-row">
      {data &&
        data.map(
          ({ title, cards }: { title: string; cards: string[] }, i: number) => {
            return (
              <div key={i}>
                <h2 className="text-xl py-4">{title}</h2>
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
    </div>
  );
}
