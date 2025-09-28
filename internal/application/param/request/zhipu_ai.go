package request

type ZhipuAiChatParam struct {
	Model       string                   `json:"model" form:"model"` //所要调用的模型编码
	Messages    []ZhipuAiChatMessageItem `json:"messages" form:"messages"`
	RequestId   string                   `json:"request_id,omitempty" form:"request_id,omitempty"`   //由用户端传参，需保证唯一性；用于区分每次请求的唯一标识，用户端不传时平台会默认生成。
	DoSample    bool                     `json:"do_sample" form:"do_sample"`                         //do_sample 为 true 时启用采样策略，do_sample 为 false 时采样策略 temperature、top_p 将不生效。默认值为 true。
	Stream      bool                     `json:"stream" form:"stream"`                               //使用同步调用时，此参数应当设置为 fasle 或者省略。表示模型生成完所有内容后一次性返回所有内容。默认值为 false。
	Temperature float32                  `json:"temperature,omitempty" form:"temperature,omitempty"` //采样温度，控制输出的随机性，必须为正数。取值范围是：[0.0, 1.0]，默认值为 0.95，值越大，会使输出更随机，更具创造性；值越小，输出会更加稳定或确定。建议您根据应用场景调整 top_p 或 temperature 参数，但不要同时调整两个参数
	TopP        float32                  `json:"top_p,omitempty" form:"top_p,omitempty"`             //用温度取样的另一种方法，称为核取样 取值范围是：[0.0, 1.0] ，默认值为 0.7。模型考虑具有 top_p 概率质量 tokens 的结果。例如：0.1 意味着模型解码器只考虑从前 10% 的概率的候选集中取 tokens
	MaxTokens   int                      `json:"max_tokens,omitempty" form:"max_tokens,omitempty"`   //模型输出最大 tokens，最大输出为4095，默认值为1024
	Stop        []string                 `json:"stop,omitempty" form:"stop,omitempty"`               //模型在遇到stop所制定的字符时将停止生成，目前仅支持单个停止词，格式为["stop_word1"]
	Tools       []ZhipuAiChatToolItem    `json:"tools,omitempty" form:"tools,omitempty"`             //模型可调用的工具列表。
	ToolChoice  any                      `json:"tool_choice,omitempty" form:"tool_choice,omitempty"` //用于控制模型是如何选择要调用的函数，仅当工具类型为function时补充。默认为auto，当前仅支持auto
	UserId      string                   `json:"user_id,omitempty" form:"user_id,omitempty"`         //终端用户的唯一ID，协助平台对终端用户的违规行为、生成违法及不良信息或其他滥用行为进行干预。ID长度要求：最少6个字符，最多128个字符
}

type ZhipuAiViewParam struct {
	Model  string `json:"model" form:"model"`                         //所要调用的模型编码
	Prompt string `json:"prompt" form:"prompt"`                       //所需图像的文本描述
	Size   string `json:"size" form:"size"`                           //图片尺寸，仅 cogview-3-plus 支持该参数。 可选范围： [1024x1024,768x1344,864x1152,1344x768,1152x864,1440x720,720x1440]，默认是1024x1024。
	UserId string `json:"user_id,omitempty" form:"user_id,omitempty"` //终端用户的唯一ID，协助平台对终端用户的违规行为、生成违法及不良信息或其他滥用行为进行干预。ID长度要求：最少6个字符，最多128个字符
}

type ZhipuAiChatMessageItem struct {
	Role    string `json:"role" form:"role"`       //角色
	Content string `json:"content" form:"content"` //聊天内容
}

type ZhipuAiChatToolItem struct {
	Type      string                    `json:"type,omitempty" form:"type,omitempty"`             //工具类型,目前支持function、retrieval、web_search
	Function  *ZhipuAiChatToolFunction  `json:"function,omitempty" form:"function,omitempty"`     //仅当工具类型为function时补充
	Retrieval *ZhipuAiChatToolRetrieval `json:"retrieval,omitempty" form:"retrieval,omitempty"`   //仅当工具类型为retrieval时补充
	WebSearch *ZhipuAiChatToolWebSearch `json:"web_search,omitempty" form:"web_search,omitempty"` //仅当工具类型为web_search时补充
}

type ZhipuAiChatToolFunction struct {
	Name        string `json:"name" form:"name"`               //函数名称，只能包含a-z，A-Z，0-9，下划线和中横线。最大长度限制为64
	Description string `json:"description" form:"description"` //用于描述函数功能。模型会根据这段描述决定函数调用方式。
	Parameters  any    `json:"parameters" form:"parameters"`   //parameter 字段需要传入一个 Json Schema 对象，以准确地定义函数所接受的参数。
}

type ZhipuAiChatToolRetrieval struct {
	KnowledgeId    string `json:"knowledge_id" form:"knowledge_id"`       //当涉及到知识库ID时，请前往开放平台的知识库模块进行创建或获取。
	PromptTemplate string `json:"prompt_template" form:"prompt_template"` //请求模型时的知识库模板
}

type ZhipuAiChatToolWebSearch struct {
	Enable       bool   `json:"enable" form:"enable"`               //网络搜索功能：默认为关闭状态（False）
	SearchQuery  string `json:"search_query" form:"search_query"`   //强制搜索自定义关键内容，此时模型会根据自定义搜索关键内容返回的结果作为背景知识来回答用户发起的对话。
	SearchResult bool   `json:"search_result" form:"search_result"` //获取详细的网页搜索来源信息，包括来源网站的图标、标题、链接、来源名称以及引用的文本内容。默认为关闭。
}
