type Success<T> = {
    data: T;
    err: null;
};

type Failure<E> = {
    data: null;
    err: E;
};

type Result<T, E = Error> = Success<T> | Failure<E>;

export async function tryCatch<T, E = Error>(
    promise: Promise<T>,
): Promise<Result<T, E>> {
    try {
        const data = await promise;
        return {data, err: null};
    } catch (err) {
        return {data: null, err: err as E};
    }
}