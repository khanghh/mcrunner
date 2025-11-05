/*eslint-disable block-scoped-var, id-length, no-control-regex, no-magic-numbers, no-prototype-builtins, no-redeclare, no-shadow, no-var, sort-vars*/
import * as $protobuf from "protobufjs/minimal";

// Common aliases
const $Reader = $protobuf.Reader, $Writer = $protobuf.Writer, $util = $protobuf.util;

// Exported root namespace
const $root = $protobuf.roots["default"] || ($protobuf.roots["default"] = {});

/**
 * MessageType enum.
 * @exports MessageType
 * @enum {number}
 * @property {number} UNKNOWN=0 UNKNOWN value
 * @property {number} ERROR=1 ERROR value
 * @property {number} PTY_BUFFER=101 PTY_BUFFER value
 * @property {number} PTY_INPUT=102 PTY_INPUT value
 * @property {number} PTY_RESIZE=103 PTY_RESIZE value
 */
export const MessageType = $root.MessageType = (() => {
    const valuesById = {}, values = Object.create(valuesById);
    values[valuesById[0] = "UNKNOWN"] = 0;
    values[valuesById[1] = "ERROR"] = 1;
    values[valuesById[101] = "PTY_BUFFER"] = 101;
    values[valuesById[102] = "PTY_INPUT"] = 102;
    values[valuesById[103] = "PTY_RESIZE"] = 103;
    return values;
})();

export const Message = $root.Message = (() => {

    /**
     * Properties of a Message.
     * @exports IMessage
     * @interface IMessage
     * @property {MessageType|null} [type] Message type
     * @property {string|null} [error] Message error
     * @property {IPtyBuffer|null} [ptyBuffer] Message ptyBuffer
     * @property {IPtyInput|null} [ptyInput] Message ptyInput
     * @property {IPtyResize|null} [ptyResize] Message ptyResize
     */

    /**
     * Constructs a new Message.
     * @exports Message
     * @classdesc Represents a Message.
     * @implements IMessage
     * @constructor
     * @param {IMessage=} [properties] Properties to set
     */
    function Message(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Message type.
     * @member {MessageType} type
     * @memberof Message
     * @instance
     */
    Message.prototype.type = 0;

    /**
     * Message error.
     * @member {string} error
     * @memberof Message
     * @instance
     */
    Message.prototype.error = "";

    /**
     * Message ptyBuffer.
     * @member {IPtyBuffer|null|undefined} ptyBuffer
     * @memberof Message
     * @instance
     */
    Message.prototype.ptyBuffer = null;

    /**
     * Message ptyInput.
     * @member {IPtyInput|null|undefined} ptyInput
     * @memberof Message
     * @instance
     */
    Message.prototype.ptyInput = null;

    /**
     * Message ptyResize.
     * @member {IPtyResize|null|undefined} ptyResize
     * @memberof Message
     * @instance
     */
    Message.prototype.ptyResize = null;

    // OneOf field names bound to virtual getters and setters
    let $oneOfFields;

    /**
     * Message payload.
     * @member {"ptyBuffer"|"ptyInput"|"ptyResize"|undefined} payload
     * @memberof Message
     * @instance
     */
    Object.defineProperty(Message.prototype, "payload", {
        get: $util.oneOfGetter($oneOfFields = ["ptyBuffer", "ptyInput", "ptyResize"]),
        set: $util.oneOfSetter($oneOfFields)
    });

    /**
     * Creates a new Message instance using the specified properties.
     * @function create
     * @memberof Message
     * @static
     * @param {IMessage=} [properties] Properties to set
     * @returns {Message} Message instance
     */
    Message.create = function create(properties) {
        return new Message(properties);
    };

    /**
     * Encodes the specified Message message. Does not implicitly {@link Message.verify|verify} messages.
     * @function encode
     * @memberof Message
     * @static
     * @param {IMessage} message Message message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Message.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.type != null && Object.hasOwnProperty.call(message, "type"))
            writer.uint32(/* id 1, wireType 0 =*/8).int32(message.type);
        if (message.error != null && Object.hasOwnProperty.call(message, "error"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.error);
        if (message.ptyBuffer != null && Object.hasOwnProperty.call(message, "ptyBuffer"))
            $root.PtyBuffer.encode(message.ptyBuffer, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
        if (message.ptyInput != null && Object.hasOwnProperty.call(message, "ptyInput"))
            $root.PtyInput.encode(message.ptyInput, writer.uint32(/* id 4, wireType 2 =*/34).fork()).ldelim();
        if (message.ptyResize != null && Object.hasOwnProperty.call(message, "ptyResize"))
            $root.PtyResize.encode(message.ptyResize, writer.uint32(/* id 5, wireType 2 =*/42).fork()).ldelim();
        return writer;
    };

    /**
     * Encodes the specified Message message, length delimited. Does not implicitly {@link Message.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Message
     * @static
     * @param {IMessage} message Message message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Message.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Message message from the specified reader or buffer.
     * @function decode
     * @memberof Message
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Message} Message
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Message.decode = function decode(reader, length, error) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.Message();
        while (reader.pos < end) {
            let tag = reader.uint32();
            if (tag === error)
                break;
            switch (tag >>> 3) {
            case 1: {
                    message.type = reader.int32();
                    break;
                }
            case 2: {
                    message.error = reader.string();
                    break;
                }
            case 3: {
                    message.ptyBuffer = $root.PtyBuffer.decode(reader, reader.uint32());
                    break;
                }
            case 4: {
                    message.ptyInput = $root.PtyInput.decode(reader, reader.uint32());
                    break;
                }
            case 5: {
                    message.ptyResize = $root.PtyResize.decode(reader, reader.uint32());
                    break;
                }
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Message message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Message
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Message} Message
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Message.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Message message.
     * @function verify
     * @memberof Message
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Message.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        let properties = {};
        if (message.type != null && message.hasOwnProperty("type"))
            switch (message.type) {
            default:
                return "type: enum value expected";
            case 0:
            case 1:
            case 101:
            case 102:
            case 103:
                break;
            }
        if (message.error != null && message.hasOwnProperty("error"))
            if (!$util.isString(message.error))
                return "error: string expected";
        if (message.ptyBuffer != null && message.hasOwnProperty("ptyBuffer")) {
            properties.payload = 1;
            {
                let error = $root.PtyBuffer.verify(message.ptyBuffer);
                if (error)
                    return "ptyBuffer." + error;
            }
        }
        if (message.ptyInput != null && message.hasOwnProperty("ptyInput")) {
            if (properties.payload === 1)
                return "payload: multiple values";
            properties.payload = 1;
            {
                let error = $root.PtyInput.verify(message.ptyInput);
                if (error)
                    return "ptyInput." + error;
            }
        }
        if (message.ptyResize != null && message.hasOwnProperty("ptyResize")) {
            if (properties.payload === 1)
                return "payload: multiple values";
            properties.payload = 1;
            {
                let error = $root.PtyResize.verify(message.ptyResize);
                if (error)
                    return "ptyResize." + error;
            }
        }
        return null;
    };

    /**
     * Creates a Message message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Message
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Message} Message
     */
    Message.fromObject = function fromObject(object) {
        if (object instanceof $root.Message)
            return object;
        let message = new $root.Message();
        switch (object.type) {
        default:
            if (typeof object.type === "number") {
                message.type = object.type;
                break;
            }
            break;
        case "UNKNOWN":
        case 0:
            message.type = 0;
            break;
        case "ERROR":
        case 1:
            message.type = 1;
            break;
        case "PTY_BUFFER":
        case 101:
            message.type = 101;
            break;
        case "PTY_INPUT":
        case 102:
            message.type = 102;
            break;
        case "PTY_RESIZE":
        case 103:
            message.type = 103;
            break;
        }
        if (object.error != null)
            message.error = String(object.error);
        if (object.ptyBuffer != null) {
            if (typeof object.ptyBuffer !== "object")
                throw TypeError(".Message.ptyBuffer: object expected");
            message.ptyBuffer = $root.PtyBuffer.fromObject(object.ptyBuffer);
        }
        if (object.ptyInput != null) {
            if (typeof object.ptyInput !== "object")
                throw TypeError(".Message.ptyInput: object expected");
            message.ptyInput = $root.PtyInput.fromObject(object.ptyInput);
        }
        if (object.ptyResize != null) {
            if (typeof object.ptyResize !== "object")
                throw TypeError(".Message.ptyResize: object expected");
            message.ptyResize = $root.PtyResize.fromObject(object.ptyResize);
        }
        return message;
    };

    /**
     * Creates a plain object from a Message message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Message
     * @static
     * @param {Message} message Message
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Message.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.type = options.enums === String ? "UNKNOWN" : 0;
            object.error = "";
        }
        if (message.type != null && message.hasOwnProperty("type"))
            object.type = options.enums === String ? $root.MessageType[message.type] === undefined ? message.type : $root.MessageType[message.type] : message.type;
        if (message.error != null && message.hasOwnProperty("error"))
            object.error = message.error;
        if (message.ptyBuffer != null && message.hasOwnProperty("ptyBuffer")) {
            object.ptyBuffer = $root.PtyBuffer.toObject(message.ptyBuffer, options);
            if (options.oneofs)
                object.payload = "ptyBuffer";
        }
        if (message.ptyInput != null && message.hasOwnProperty("ptyInput")) {
            object.ptyInput = $root.PtyInput.toObject(message.ptyInput, options);
            if (options.oneofs)
                object.payload = "ptyInput";
        }
        if (message.ptyResize != null && message.hasOwnProperty("ptyResize")) {
            object.ptyResize = $root.PtyResize.toObject(message.ptyResize, options);
            if (options.oneofs)
                object.payload = "ptyResize";
        }
        return object;
    };

    /**
     * Converts this Message to JSON.
     * @function toJSON
     * @memberof Message
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Message.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    /**
     * Gets the default type url for Message
     * @function getTypeUrl
     * @memberof Message
     * @static
     * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
     * @returns {string} The default type url
     */
    Message.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
        if (typeUrlPrefix === undefined) {
            typeUrlPrefix = "type.googleapis.com";
        }
        return typeUrlPrefix + "/Message";
    };

    return Message;
})();

export const PtyBuffer = $root.PtyBuffer = (() => {

    /**
     * Properties of a PtyBuffer.
     * @exports IPtyBuffer
     * @interface IPtyBuffer
     * @property {string|null} [sessionId] PtyBuffer sessionId
     * @property {Uint8Array|null} [data] PtyBuffer data
     */

    /**
     * Constructs a new PtyBuffer.
     * @exports PtyBuffer
     * @classdesc Represents a PtyBuffer.
     * @implements IPtyBuffer
     * @constructor
     * @param {IPtyBuffer=} [properties] Properties to set
     */
    function PtyBuffer(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * PtyBuffer sessionId.
     * @member {string} sessionId
     * @memberof PtyBuffer
     * @instance
     */
    PtyBuffer.prototype.sessionId = "";

    /**
     * PtyBuffer data.
     * @member {Uint8Array} data
     * @memberof PtyBuffer
     * @instance
     */
    PtyBuffer.prototype.data = $util.newBuffer([]);

    /**
     * Creates a new PtyBuffer instance using the specified properties.
     * @function create
     * @memberof PtyBuffer
     * @static
     * @param {IPtyBuffer=} [properties] Properties to set
     * @returns {PtyBuffer} PtyBuffer instance
     */
    PtyBuffer.create = function create(properties) {
        return new PtyBuffer(properties);
    };

    /**
     * Encodes the specified PtyBuffer message. Does not implicitly {@link PtyBuffer.verify|verify} messages.
     * @function encode
     * @memberof PtyBuffer
     * @static
     * @param {IPtyBuffer} message PtyBuffer message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    PtyBuffer.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.sessionId != null && Object.hasOwnProperty.call(message, "sessionId"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.sessionId);
        if (message.data != null && Object.hasOwnProperty.call(message, "data"))
            writer.uint32(/* id 2, wireType 2 =*/18).bytes(message.data);
        return writer;
    };

    /**
     * Encodes the specified PtyBuffer message, length delimited. Does not implicitly {@link PtyBuffer.verify|verify} messages.
     * @function encodeDelimited
     * @memberof PtyBuffer
     * @static
     * @param {IPtyBuffer} message PtyBuffer message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    PtyBuffer.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a PtyBuffer message from the specified reader or buffer.
     * @function decode
     * @memberof PtyBuffer
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {PtyBuffer} PtyBuffer
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    PtyBuffer.decode = function decode(reader, length, error) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.PtyBuffer();
        while (reader.pos < end) {
            let tag = reader.uint32();
            if (tag === error)
                break;
            switch (tag >>> 3) {
            case 1: {
                    message.sessionId = reader.string();
                    break;
                }
            case 2: {
                    message.data = reader.bytes();
                    break;
                }
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a PtyBuffer message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof PtyBuffer
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {PtyBuffer} PtyBuffer
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    PtyBuffer.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a PtyBuffer message.
     * @function verify
     * @memberof PtyBuffer
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    PtyBuffer.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.sessionId != null && message.hasOwnProperty("sessionId"))
            if (!$util.isString(message.sessionId))
                return "sessionId: string expected";
        if (message.data != null && message.hasOwnProperty("data"))
            if (!(message.data && typeof message.data.length === "number" || $util.isString(message.data)))
                return "data: buffer expected";
        return null;
    };

    /**
     * Creates a PtyBuffer message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof PtyBuffer
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {PtyBuffer} PtyBuffer
     */
    PtyBuffer.fromObject = function fromObject(object) {
        if (object instanceof $root.PtyBuffer)
            return object;
        let message = new $root.PtyBuffer();
        if (object.sessionId != null)
            message.sessionId = String(object.sessionId);
        if (object.data != null)
            if (typeof object.data === "string")
                $util.base64.decode(object.data, message.data = $util.newBuffer($util.base64.length(object.data)), 0);
            else if (object.data.length >= 0)
                message.data = object.data;
        return message;
    };

    /**
     * Creates a plain object from a PtyBuffer message. Also converts values to other types if specified.
     * @function toObject
     * @memberof PtyBuffer
     * @static
     * @param {PtyBuffer} message PtyBuffer
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    PtyBuffer.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.sessionId = "";
            if (options.bytes === String)
                object.data = "";
            else {
                object.data = [];
                if (options.bytes !== Array)
                    object.data = $util.newBuffer(object.data);
            }
        }
        if (message.sessionId != null && message.hasOwnProperty("sessionId"))
            object.sessionId = message.sessionId;
        if (message.data != null && message.hasOwnProperty("data"))
            object.data = options.bytes === String ? $util.base64.encode(message.data, 0, message.data.length) : options.bytes === Array ? Array.prototype.slice.call(message.data) : message.data;
        return object;
    };

    /**
     * Converts this PtyBuffer to JSON.
     * @function toJSON
     * @memberof PtyBuffer
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    PtyBuffer.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    /**
     * Gets the default type url for PtyBuffer
     * @function getTypeUrl
     * @memberof PtyBuffer
     * @static
     * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
     * @returns {string} The default type url
     */
    PtyBuffer.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
        if (typeUrlPrefix === undefined) {
            typeUrlPrefix = "type.googleapis.com";
        }
        return typeUrlPrefix + "/PtyBuffer";
    };

    return PtyBuffer;
})();

export const PtyInput = $root.PtyInput = (() => {

    /**
     * Properties of a PtyInput.
     * @exports IPtyInput
     * @interface IPtyInput
     * @property {string|null} [sessionId] PtyInput sessionId
     * @property {Uint8Array|null} [data] PtyInput data
     */

    /**
     * Constructs a new PtyInput.
     * @exports PtyInput
     * @classdesc Represents a PtyInput.
     * @implements IPtyInput
     * @constructor
     * @param {IPtyInput=} [properties] Properties to set
     */
    function PtyInput(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * PtyInput sessionId.
     * @member {string} sessionId
     * @memberof PtyInput
     * @instance
     */
    PtyInput.prototype.sessionId = "";

    /**
     * PtyInput data.
     * @member {Uint8Array} data
     * @memberof PtyInput
     * @instance
     */
    PtyInput.prototype.data = $util.newBuffer([]);

    /**
     * Creates a new PtyInput instance using the specified properties.
     * @function create
     * @memberof PtyInput
     * @static
     * @param {IPtyInput=} [properties] Properties to set
     * @returns {PtyInput} PtyInput instance
     */
    PtyInput.create = function create(properties) {
        return new PtyInput(properties);
    };

    /**
     * Encodes the specified PtyInput message. Does not implicitly {@link PtyInput.verify|verify} messages.
     * @function encode
     * @memberof PtyInput
     * @static
     * @param {IPtyInput} message PtyInput message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    PtyInput.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.sessionId != null && Object.hasOwnProperty.call(message, "sessionId"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.sessionId);
        if (message.data != null && Object.hasOwnProperty.call(message, "data"))
            writer.uint32(/* id 2, wireType 2 =*/18).bytes(message.data);
        return writer;
    };

    /**
     * Encodes the specified PtyInput message, length delimited. Does not implicitly {@link PtyInput.verify|verify} messages.
     * @function encodeDelimited
     * @memberof PtyInput
     * @static
     * @param {IPtyInput} message PtyInput message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    PtyInput.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a PtyInput message from the specified reader or buffer.
     * @function decode
     * @memberof PtyInput
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {PtyInput} PtyInput
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    PtyInput.decode = function decode(reader, length, error) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.PtyInput();
        while (reader.pos < end) {
            let tag = reader.uint32();
            if (tag === error)
                break;
            switch (tag >>> 3) {
            case 1: {
                    message.sessionId = reader.string();
                    break;
                }
            case 2: {
                    message.data = reader.bytes();
                    break;
                }
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a PtyInput message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof PtyInput
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {PtyInput} PtyInput
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    PtyInput.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a PtyInput message.
     * @function verify
     * @memberof PtyInput
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    PtyInput.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.sessionId != null && message.hasOwnProperty("sessionId"))
            if (!$util.isString(message.sessionId))
                return "sessionId: string expected";
        if (message.data != null && message.hasOwnProperty("data"))
            if (!(message.data && typeof message.data.length === "number" || $util.isString(message.data)))
                return "data: buffer expected";
        return null;
    };

    /**
     * Creates a PtyInput message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof PtyInput
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {PtyInput} PtyInput
     */
    PtyInput.fromObject = function fromObject(object) {
        if (object instanceof $root.PtyInput)
            return object;
        let message = new $root.PtyInput();
        if (object.sessionId != null)
            message.sessionId = String(object.sessionId);
        if (object.data != null)
            if (typeof object.data === "string")
                $util.base64.decode(object.data, message.data = $util.newBuffer($util.base64.length(object.data)), 0);
            else if (object.data.length >= 0)
                message.data = object.data;
        return message;
    };

    /**
     * Creates a plain object from a PtyInput message. Also converts values to other types if specified.
     * @function toObject
     * @memberof PtyInput
     * @static
     * @param {PtyInput} message PtyInput
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    PtyInput.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.sessionId = "";
            if (options.bytes === String)
                object.data = "";
            else {
                object.data = [];
                if (options.bytes !== Array)
                    object.data = $util.newBuffer(object.data);
            }
        }
        if (message.sessionId != null && message.hasOwnProperty("sessionId"))
            object.sessionId = message.sessionId;
        if (message.data != null && message.hasOwnProperty("data"))
            object.data = options.bytes === String ? $util.base64.encode(message.data, 0, message.data.length) : options.bytes === Array ? Array.prototype.slice.call(message.data) : message.data;
        return object;
    };

    /**
     * Converts this PtyInput to JSON.
     * @function toJSON
     * @memberof PtyInput
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    PtyInput.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    /**
     * Gets the default type url for PtyInput
     * @function getTypeUrl
     * @memberof PtyInput
     * @static
     * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
     * @returns {string} The default type url
     */
    PtyInput.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
        if (typeUrlPrefix === undefined) {
            typeUrlPrefix = "type.googleapis.com";
        }
        return typeUrlPrefix + "/PtyInput";
    };

    return PtyInput;
})();

export const PtyResize = $root.PtyResize = (() => {

    /**
     * Properties of a PtyResize.
     * @exports IPtyResize
     * @interface IPtyResize
     * @property {string|null} [sessionId] PtyResize sessionId
     * @property {number|null} [cols] PtyResize cols
     * @property {number|null} [rows] PtyResize rows
     */

    /**
     * Constructs a new PtyResize.
     * @exports PtyResize
     * @classdesc Represents a PtyResize.
     * @implements IPtyResize
     * @constructor
     * @param {IPtyResize=} [properties] Properties to set
     */
    function PtyResize(properties) {
        if (properties)
            for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * PtyResize sessionId.
     * @member {string} sessionId
     * @memberof PtyResize
     * @instance
     */
    PtyResize.prototype.sessionId = "";

    /**
     * PtyResize cols.
     * @member {number} cols
     * @memberof PtyResize
     * @instance
     */
    PtyResize.prototype.cols = 0;

    /**
     * PtyResize rows.
     * @member {number} rows
     * @memberof PtyResize
     * @instance
     */
    PtyResize.prototype.rows = 0;

    /**
     * Creates a new PtyResize instance using the specified properties.
     * @function create
     * @memberof PtyResize
     * @static
     * @param {IPtyResize=} [properties] Properties to set
     * @returns {PtyResize} PtyResize instance
     */
    PtyResize.create = function create(properties) {
        return new PtyResize(properties);
    };

    /**
     * Encodes the specified PtyResize message. Does not implicitly {@link PtyResize.verify|verify} messages.
     * @function encode
     * @memberof PtyResize
     * @static
     * @param {IPtyResize} message PtyResize message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    PtyResize.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.sessionId != null && Object.hasOwnProperty.call(message, "sessionId"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.sessionId);
        if (message.cols != null && Object.hasOwnProperty.call(message, "cols"))
            writer.uint32(/* id 2, wireType 0 =*/16).uint32(message.cols);
        if (message.rows != null && Object.hasOwnProperty.call(message, "rows"))
            writer.uint32(/* id 3, wireType 0 =*/24).uint32(message.rows);
        return writer;
    };

    /**
     * Encodes the specified PtyResize message, length delimited. Does not implicitly {@link PtyResize.verify|verify} messages.
     * @function encodeDelimited
     * @memberof PtyResize
     * @static
     * @param {IPtyResize} message PtyResize message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    PtyResize.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a PtyResize message from the specified reader or buffer.
     * @function decode
     * @memberof PtyResize
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {PtyResize} PtyResize
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    PtyResize.decode = function decode(reader, length, error) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        let end = length === undefined ? reader.len : reader.pos + length, message = new $root.PtyResize();
        while (reader.pos < end) {
            let tag = reader.uint32();
            if (tag === error)
                break;
            switch (tag >>> 3) {
            case 1: {
                    message.sessionId = reader.string();
                    break;
                }
            case 2: {
                    message.cols = reader.uint32();
                    break;
                }
            case 3: {
                    message.rows = reader.uint32();
                    break;
                }
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a PtyResize message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof PtyResize
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {PtyResize} PtyResize
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    PtyResize.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a PtyResize message.
     * @function verify
     * @memberof PtyResize
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    PtyResize.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.sessionId != null && message.hasOwnProperty("sessionId"))
            if (!$util.isString(message.sessionId))
                return "sessionId: string expected";
        if (message.cols != null && message.hasOwnProperty("cols"))
            if (!$util.isInteger(message.cols))
                return "cols: integer expected";
        if (message.rows != null && message.hasOwnProperty("rows"))
            if (!$util.isInteger(message.rows))
                return "rows: integer expected";
        return null;
    };

    /**
     * Creates a PtyResize message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof PtyResize
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {PtyResize} PtyResize
     */
    PtyResize.fromObject = function fromObject(object) {
        if (object instanceof $root.PtyResize)
            return object;
        let message = new $root.PtyResize();
        if (object.sessionId != null)
            message.sessionId = String(object.sessionId);
        if (object.cols != null)
            message.cols = object.cols >>> 0;
        if (object.rows != null)
            message.rows = object.rows >>> 0;
        return message;
    };

    /**
     * Creates a plain object from a PtyResize message. Also converts values to other types if specified.
     * @function toObject
     * @memberof PtyResize
     * @static
     * @param {PtyResize} message PtyResize
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    PtyResize.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        let object = {};
        if (options.defaults) {
            object.sessionId = "";
            object.cols = 0;
            object.rows = 0;
        }
        if (message.sessionId != null && message.hasOwnProperty("sessionId"))
            object.sessionId = message.sessionId;
        if (message.cols != null && message.hasOwnProperty("cols"))
            object.cols = message.cols;
        if (message.rows != null && message.hasOwnProperty("rows"))
            object.rows = message.rows;
        return object;
    };

    /**
     * Converts this PtyResize to JSON.
     * @function toJSON
     * @memberof PtyResize
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    PtyResize.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    /**
     * Gets the default type url for PtyResize
     * @function getTypeUrl
     * @memberof PtyResize
     * @static
     * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
     * @returns {string} The default type url
     */
    PtyResize.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
        if (typeUrlPrefix === undefined) {
            typeUrlPrefix = "type.googleapis.com";
        }
        return typeUrlPrefix + "/PtyResize";
    };

    return PtyResize;
})();

export { $root as default };
