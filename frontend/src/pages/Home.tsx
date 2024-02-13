// Dependencies
import { useQuery } from "@tanstack/react-query";
import axios from "axios";

export const Home = () => {
  return (
    <>
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
    <div className="flex w-full flex-row justify-evenly">
      {data &&
        data.map(
          ({ title, cards }: { title: string; cards: string[] }, i: number) => {
            return (
              <div key={i}>
                <h2 className="py-4 text-xl">{title}</h2>
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
          },
        )}
    </div>
  );
}
