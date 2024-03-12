// Dependencies
import { useMutation } from "@tanstack/react-query";
import axios from "axios";

// Component
import { Button, IButtonStyles } from "../components/Button";

export const Home = () => {
  return (
    <>
      <CardUpgrades />
    </>
  );
};

function CardUpgrades(): JSX.Element {
  const { error, data } = useMutation({
    mutationFn: async ({ data }: any) => {
      return axios
        .post("http://localhost:8080/api/cards/upgrades", data, {
          headers: {
            "Content-Type": "application/json",
          },
        })
        .then(({ data }) => {
          return data;
        });
    },
  });

  if (error) {
    return <div>{error.toString()}</div>;
  }

  /**
   * TODO:
   * 1. Create a form with two inputs: commander and deck name
   * 2. Create a button to submit the form
   * 3. Create a function to handle the form submission
   * 4. Separate the form from the component that displays the card upgrades
   */

  return (
    <main>
      <div className="max-w-lg flex-row">
        <form className="flex flex-row justify-center items-center bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
          <div className="mx-5">
            <input
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              id="commander"
              type="text"
              placeholder="e.g Kardur, Doomscourge"
            />
          </div>
          <div className="mx-5">
            <input
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              id="deck-name"
              type="text"
              placeholder="e.g Chaos Incarnate"
            />
          </div>
          <Button
            text="Submit"
            styles="cursor-pointer rounded-lg border bg-white px-3 tracking-wide h-10"
            onClick={() => {}}
          />
        </form>
      </div>
      <div className="flex w-full flex-row justify-evenly">
        {data &&
          data.map(
            (
              { title, cards }: { title: string; cards: string[] },
              i: number
            ) => {
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
            }
          )}
      </div>
    </main>
  );
}
