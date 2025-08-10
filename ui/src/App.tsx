import { Menubar } from "@/components/ui/menubar";
import { Button } from "./components/ui/button";
import { Input } from "./components/ui/input";
import * as Comlink from "comlink"
import type CardQuery from "./workers/cardQuery";
import react from "react";

function CardSearch({ cardQuery, setCardUrl }: { cardQuery: Comlink.Remote<CardQuery>, setCardUrl: react.Dispatch<react.SetStateAction<string | null>> }) {
  const searchRef = react.useRef<HTMLInputElement>(null)
  async function handleSubmit(e: react.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const query = searchRef.current?.value
    if (query == null) {
      alert("query null")
      return
    }
    const cards = await cardQuery.queryCards(query)
    if (cards.length > 0) {
      setCardUrl(await cardQuery.getCard(cards[0]).then(c => c.image_uris.normal))
    } else {
      setCardUrl(null)
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

function TitleBar({ cardQuery, setCardUrl }: { cardQuery: Comlink.Remote<CardQuery>, setCardUrl: react.Dispatch<react.SetStateAction<string | null>> }) {
  return (
    <>
      <Menubar>
        <Button asChild>
          <a href="/">MtgBuilder</a>
        </Button>
        <CardSearch cardQuery={cardQuery} setCardUrl={setCardUrl} />
      </Menubar >
    </>
  );
}

function App({ cardQuery }: { cardQuery: Comlink.Remote<CardQuery> }) {
  const [cardUrl, setCardUrl] = react.useState<string | null>(null)
  return (
    <>
      <TitleBar cardQuery={cardQuery} setCardUrl={setCardUrl} />
      {
        cardUrl != null && <img src={cardUrl} />
      }
    </>
  );
}

export default App;
