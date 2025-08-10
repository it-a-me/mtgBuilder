import { Menubar } from "@/components/ui/menubar";
import { Button } from "./components/ui/button";
import { Input } from "./components/ui/input";
import * as Comlink from "comlink"
import type CardQuery from "./workers/cardQuery";
import react from "react";

function CardSearch({ cardQuery, setCardUrls }: { cardQuery: Comlink.Remote<CardQuery>, setCardUrls: react.Dispatch<react.SetStateAction<string[]>> }) {
  const searchRef = react.useRef<HTMLInputElement>(null)
  async function handleSubmit(e: react.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const query = searchRef.current?.value
    if (query == null) {
      alert("query null")
      return
    }
    const cards = await cardQuery.queryCards(query)
    cards.length = Math.min(cards.length, 30)
    setCardUrls([])
    for (const c of cards) {
      const card = await cardQuery.getCard(c)
      const url = card.image_uris.normal
      setCardUrls(u => u.concat([url]))
    }
  }
  return (
    <>
      <form className="flex w-full" onSubmit={e => handleSubmit(e)}>
        <Input ref={searchRef} type="search" placeholder="name:goblin type:creature oracle:/create.*token/" />
        <Button type="submit" variant="outline" >
          Search
        </Button>
      </form>
    </>
  );
}

function TitleBar({ cardQuery, setCardUrls }: { cardQuery: Comlink.Remote<CardQuery>, setCardUrls: react.Dispatch<react.SetStateAction<string[]>> }) {
  return (
    <>
      <Menubar className="pl-0">
        <Button className="bg-blue-300 hover:bg-blue-400 text-gray-800 font-bold" asChild>
          <a href="/">MtgBuilder</a>
        </Button>
        <CardSearch cardQuery={cardQuery} setCardUrls={setCardUrls} />
      </Menubar >
    </>
  );
}

function App({ cardQuery }: { cardQuery: Comlink.Remote<CardQuery> }) {
  const [cardUrls, setCardUrls] = react.useState<string[]>([])
  return (
    <>
      <div className="bg-gray-800">
        <div className="pb-4">
          <TitleBar cardQuery={cardQuery} setCardUrls={setCardUrls} />
        </div>
        <div className="grid grid-cols-3">
          {
            cardUrls.map(url => <img className="w-80" src={url} />)
          }
        </div>
      </div>
    </>
  );
}

export default App;
