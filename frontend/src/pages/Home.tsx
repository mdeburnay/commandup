// Dependencies
import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import axios from "axios";

// Component
import { Button } from "../components/Buttons/PrimaryButton";

// Icons
import { AiOutlineLoading3Quarters } from "react-icons/ai";

interface ICard {
  name: string;
  synergy: number;
  inclusion: number;
}
interface ICardCategory {
  title: string;
  cards: ICard[];
}

interface GetCardUpgradesReq {
  commander: string;
  precon?: string;
}

export const Home = () => {
  return (
    <>
      <CardUpgrades />
    </>
  );
};

function CardUpgrades(): JSX.Element {
  const [commander, setCommander] = useState("");
  const [deckName, setDeckName] = useState("");
  const [error, setError] = useState("");

  const {
    mutate: fetchCardUpgrades,
    isError,
    isPending,
    data,
  } = useMutation({
    mutationFn: async (payload: GetCardUpgradesReq) => {
      return axios
        .post("http://localhost:8080/api/cards/upgrades", payload, {
          headers: {
            "Content-Type": "application/json",
          },
        })
        .then(({ data }) => {
          return data;
        });
    },
  });

  const handleSubmit = (e: any) => {
    e.preventDefault();
    const payload = {
      commander: commander,
      precon: deckName,
    };

    fetchCardUpgrades(payload);
  };

  /**
   * TODO:
   * 1. Create a form with two inputs: commander and deck name
   * 2. Create a button to submit the form
   * 3. Create a function to handle the form submission
   * 4. Separate the form from the component that displays the card upgrades
   */

  return (
    <main>
      <div className="max-w-lg flex-row bg-white shadow-md">
        <form className="flex flex-row justify-center items-cente rounded px-8 pt-6 pb-8 mb-4">
          <div className="mx-5">
            <input
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              id="commander"
              type="text"
              onChange={(e) => setCommander(e.target.value)}
              placeholder="e.g Kardur, Doomscourge"
            />
          </div>
          <div className="mx-5">
            <input
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              id="deck-name"
              type="text"
              onChange={(e) => setDeckName(e.target.value)}
              placeholder="e.g Chaos Incarnate"
            />
          </div>
          <Button
            text="Submit"
            styles="cursor-pointer rounded-lg border bg-white px-3 tracking-wide h-10"
            onClick={handleSubmit}
          />
        </form>
        <div className="flex justify-center">
          {isPending && <AiOutlineLoading3Quarters className="animate-spin" />}
          {isError && <div>Not found!</div>}
        </div>
      </div>
      <div className="flex w-full flex-row justify-evenly">
        {data &&
          data.map(({ title, cards }: ICardCategory, i: number) => {
            return (
              <div key={i}>
                <h2 className="py-4 text-xl">{title}</h2>
                <div>
                  {cards.map(
                    ({ name, synergy, inclusion }: ICard, i: number) => {
                      return (
                        <div key={i} style={{ fontSize: 14 }}>
                          {name}
                          {synergy ? `- Synergy: ${synergy}%` : ""}
                          {inclusion ? ` - Inclusion: ${inclusion}` : ""}
                        </div>
                      );
                    }
                  )}
                </div>
              </div>
            );
          })}
      </div>
    </main>
  );
}
