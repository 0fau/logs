import {formatDistance} from "date-fns";
import numeral from "numeral";

export function formatDate(date: number): string {
    return formatDistance(new Date(date), new Date(), {addSuffix: true})
}

export function formatDuration(duration: number): string {
    let date = new Date(duration);
    return date.getMinutes() + 'm' + date.getSeconds() + 's';
}

export function formatSeconds(milli: number): string {
    return numeral(milli / 1000).format('0.0a')
}

export function formatDamage(damage: number): string {
    return numeral(damage).format('0.0a')
}

export function formatPercent(percent: number): string {
    const ret = numeral(percent * 100).format('0.0');
    return ret == "0.0" ? "" : ret;
}

export function formatPercentFlat(percent: number): string {
    const ret = numeral(percent * 100).format('0');
    return ret == "0" ? "" : ret;
}