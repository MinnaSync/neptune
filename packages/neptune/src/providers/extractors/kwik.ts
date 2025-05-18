import { errAsync } from "neverthrow";
import { evalSync } from "../../util/vm";

/**
 * Extract the m3u8 from the kwik player.
 * @param url The kwik embed url.
 */
async function extract(url: string) {
    const res = await fetch(url, {
        headers: { referer: "https://animepahe.ru/" }
    });

    if (res.status !== 200) {
        return errAsync(new Error('Failed to fetch from animepahe.'));
    }

    const html = await res.text();
    const script = html.match(/(eval)(\(f.*?)(\n<\/script>)/s)![2].replace('eval', '');

    /**
     * Dangerously evaling code is stupid, regardless of where it is from.
     * Bun doesn't support implementations like vm2 or isolated-vm.
     * This isn't the best solution either, realistically you'd want to sandbox it into another container.
     * @see https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/eval#never_use_direct_eval!
     */
    const result = evalSync(script).match(/https.*?m3u8/)[0];

    return result;
}

export default {
    extract
}