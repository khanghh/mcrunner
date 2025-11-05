import * as $protobuf from "protobufjs";
import Long = require("long");
/** MessageType enum. */
export enum MessageType {
    UNKNOWN = 0,
    ERROR = 1,
    PTY_BUFFER = 101,
    PTY_INPUT = 102,
    PTY_RESIZE = 103
}

/** Represents a Message. */
export class Message implements IMessage {

    /**
     * Constructs a new Message.
     * @param [properties] Properties to set
     */
    constructor(properties?: IMessage);

    /** Message type. */
    public type: MessageType;

    /** Message error. */
    public error: string;

    /** Message ptyBuffer. */
    public ptyBuffer?: (IPtyBuffer|null);

    /** Message ptyInput. */
    public ptyInput?: (IPtyInput|null);

    /** Message ptyResize. */
    public ptyResize?: (IPtyResize|null);

    /** Message payload. */
    public payload?: ("ptyBuffer"|"ptyInput"|"ptyResize");

    /**
     * Creates a new Message instance using the specified properties.
     * @param [properties] Properties to set
     * @returns Message instance
     */
    public static create(properties?: IMessage): Message;

    /**
     * Encodes the specified Message message. Does not implicitly {@link Message.verify|verify} messages.
     * @param message Message message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IMessage, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified Message message, length delimited. Does not implicitly {@link Message.verify|verify} messages.
     * @param message Message message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IMessage, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a Message message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns Message
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): Message;

    /**
     * Decodes a Message message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns Message
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): Message;

    /**
     * Verifies a Message message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a Message message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns Message
     */
    public static fromObject(object: { [k: string]: any }): Message;

    /**
     * Creates a plain object from a Message message. Also converts values to other types if specified.
     * @param message Message
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: Message, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this Message to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };

    /**
     * Gets the default type url for Message
     * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
     * @returns The default type url
     */
    public static getTypeUrl(typeUrlPrefix?: string): string;
}

/** Represents a PtyBuffer. */
export class PtyBuffer implements IPtyBuffer {

    /**
     * Constructs a new PtyBuffer.
     * @param [properties] Properties to set
     */
    constructor(properties?: IPtyBuffer);

    /** PtyBuffer sessionId. */
    public sessionId: string;

    /** PtyBuffer data. */
    public data: Uint8Array;

    /**
     * Creates a new PtyBuffer instance using the specified properties.
     * @param [properties] Properties to set
     * @returns PtyBuffer instance
     */
    public static create(properties?: IPtyBuffer): PtyBuffer;

    /**
     * Encodes the specified PtyBuffer message. Does not implicitly {@link PtyBuffer.verify|verify} messages.
     * @param message PtyBuffer message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IPtyBuffer, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified PtyBuffer message, length delimited. Does not implicitly {@link PtyBuffer.verify|verify} messages.
     * @param message PtyBuffer message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IPtyBuffer, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a PtyBuffer message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns PtyBuffer
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): PtyBuffer;

    /**
     * Decodes a PtyBuffer message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns PtyBuffer
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): PtyBuffer;

    /**
     * Verifies a PtyBuffer message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a PtyBuffer message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns PtyBuffer
     */
    public static fromObject(object: { [k: string]: any }): PtyBuffer;

    /**
     * Creates a plain object from a PtyBuffer message. Also converts values to other types if specified.
     * @param message PtyBuffer
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: PtyBuffer, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this PtyBuffer to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };

    /**
     * Gets the default type url for PtyBuffer
     * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
     * @returns The default type url
     */
    public static getTypeUrl(typeUrlPrefix?: string): string;
}

/** Represents a PtyInput. */
export class PtyInput implements IPtyInput {

    /**
     * Constructs a new PtyInput.
     * @param [properties] Properties to set
     */
    constructor(properties?: IPtyInput);

    /** PtyInput sessionId. */
    public sessionId: string;

    /** PtyInput data. */
    public data: Uint8Array;

    /**
     * Creates a new PtyInput instance using the specified properties.
     * @param [properties] Properties to set
     * @returns PtyInput instance
     */
    public static create(properties?: IPtyInput): PtyInput;

    /**
     * Encodes the specified PtyInput message. Does not implicitly {@link PtyInput.verify|verify} messages.
     * @param message PtyInput message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IPtyInput, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified PtyInput message, length delimited. Does not implicitly {@link PtyInput.verify|verify} messages.
     * @param message PtyInput message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IPtyInput, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a PtyInput message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns PtyInput
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): PtyInput;

    /**
     * Decodes a PtyInput message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns PtyInput
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): PtyInput;

    /**
     * Verifies a PtyInput message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a PtyInput message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns PtyInput
     */
    public static fromObject(object: { [k: string]: any }): PtyInput;

    /**
     * Creates a plain object from a PtyInput message. Also converts values to other types if specified.
     * @param message PtyInput
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: PtyInput, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this PtyInput to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };

    /**
     * Gets the default type url for PtyInput
     * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
     * @returns The default type url
     */
    public static getTypeUrl(typeUrlPrefix?: string): string;
}

/** Represents a PtyResize. */
export class PtyResize implements IPtyResize {

    /**
     * Constructs a new PtyResize.
     * @param [properties] Properties to set
     */
    constructor(properties?: IPtyResize);

    /** PtyResize sessionId. */
    public sessionId: string;

    /** PtyResize cols. */
    public cols: number;

    /** PtyResize rows. */
    public rows: number;

    /**
     * Creates a new PtyResize instance using the specified properties.
     * @param [properties] Properties to set
     * @returns PtyResize instance
     */
    public static create(properties?: IPtyResize): PtyResize;

    /**
     * Encodes the specified PtyResize message. Does not implicitly {@link PtyResize.verify|verify} messages.
     * @param message PtyResize message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encode(message: IPtyResize, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Encodes the specified PtyResize message, length delimited. Does not implicitly {@link PtyResize.verify|verify} messages.
     * @param message PtyResize message or plain object to encode
     * @param [writer] Writer to encode to
     * @returns Writer
     */
    public static encodeDelimited(message: IPtyResize, writer?: $protobuf.Writer): $protobuf.Writer;

    /**
     * Decodes a PtyResize message from the specified reader or buffer.
     * @param reader Reader or buffer to decode from
     * @param [length] Message length if known beforehand
     * @returns PtyResize
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): PtyResize;

    /**
     * Decodes a PtyResize message from the specified reader or buffer, length delimited.
     * @param reader Reader or buffer to decode from
     * @returns PtyResize
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): PtyResize;

    /**
     * Verifies a PtyResize message.
     * @param message Plain object to verify
     * @returns `null` if valid, otherwise the reason why it is not
     */
    public static verify(message: { [k: string]: any }): (string|null);

    /**
     * Creates a PtyResize message from a plain object. Also converts values to their respective internal types.
     * @param object Plain object
     * @returns PtyResize
     */
    public static fromObject(object: { [k: string]: any }): PtyResize;

    /**
     * Creates a plain object from a PtyResize message. Also converts values to other types if specified.
     * @param message PtyResize
     * @param [options] Conversion options
     * @returns Plain object
     */
    public static toObject(message: PtyResize, options?: $protobuf.IConversionOptions): { [k: string]: any };

    /**
     * Converts this PtyResize to JSON.
     * @returns JSON object
     */
    public toJSON(): { [k: string]: any };

    /**
     * Gets the default type url for PtyResize
     * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
     * @returns The default type url
     */
    public static getTypeUrl(typeUrlPrefix?: string): string;
}
