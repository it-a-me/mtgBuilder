import { Menubar } from "@/components/ui/menubar";
import { Button } from "./components/ui/button";
import { Input } from "./components/ui/input";
import react from "react";
import type CardQuery from "./cardQuery";

function CardSearch({ cardQuery, setCard: setCardUrl }: { cardQuery: CardQuery, setCard: react.Dispatch<react.SetStateAction<null>> }) {
  const searchRef = react.useRef<HTMLInputElement>(null)
  function handleSubmit(e: react.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const query = searchRef.current?.value
    const start = Date.now()
    const parsed = cardQuery.queryCards(query ? query : "")
    console.log(`queried cards in ${ (Date.now() - start).toString()}ms`)
    if (parsed.length > 0) {
      const first = cardQuery.getCard(parsed[0])
      const url = first.image_uris?.normal;
      setCardUrl(url)
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

function TitleBar({ cardQuery, setCard: setCardUrl }: { cardQuery: CardQuery, setCard: react.Dispatch<react.SetStateAction<null>> }) {
  return (
    <>
      <Menubar>
        <Button asChild>
          <a href="/">MtgBuilder</a>
        </Button>
        <CardSearch cardQuery={cardQuery} setCard={setCardUrl} />
      </Menubar >
    </>
  );
}

function App({ cardQuery }: { cardQuery: CardQuery }) {
  const [cardUrl, setCardUrl] = react.useState(null)
  return (
    <>
      <TitleBar cardQuery={cardQuery} setCard={setCardUrl} />
      {
        cardUrl != null && <img src={cardUrl} />
      }
    </>
  );
}

export default App;
