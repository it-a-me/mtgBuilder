import "./wasm_exec"
export default class CardQuery {
  constructor() {
    //@ts-expect-error untyped globalThis
    if (typeof globalThis.GO_cardQuery != 'undefined') {
      return;
    }
    //@ts-expect-error untyped wasm_exec
    const go = new Go();
    WebAssembly.instantiateStreaming(fetch("/wasmlib.wasm"), go.importObject).then((result) => {
      go.run(result.instance);
    });
  }

  feedCards(cardsJson: string) {
    // @ts-expect-error untyped value
    const res = globalThis.GO_cardQuery.feedCards(cardsJson)
    if (res instanceof (Error)) {
      throw res
    }
    if (res == null) {
      return
    }
    throw "unreachable";
  }

  queryCards(query: string): number[] {
    // @ts-expect-error untyped value
    const res = globalThis.GO_cardQuery.queryCards(query)
    if (res instanceof (Error)) {
      throw res
    }
    if (Array.isArray(res)) {
      return res
    }
    throw "unreachable";
  }

  parseQuery(query: string) {
    // @ts-expect-error untyped value
    const res = globalThis.GO_cardQuery.parseQuery(query)
    if (res instanceof (Error)) {
      throw res
    }
    if (typeof res == 'string') {
      return JSON.parse(res);
    }
    throw "unreachable";
  }

  getCard(cardIndex: number) {
    // @ts-expect-error untyped value
    const res = globalThis.GO_cardQuery.getCard(cardIndex)
    if (res instanceof (Error)) {
      throw res
    }
    if (typeof res == 'string') {
      return JSON.parse(res);
    }
    throw new Error("unimplemented");
  }
}

