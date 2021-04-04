<template>
  <div class="home">
    <b-container>
      <b-card class="mt-2">
        <b-form-textarea
          ref="content"
          id="content"
          v-model="content"
          rows="10"
          max-rows="10"
          plaintext
        ></b-form-textarea>
      </b-card>
      <b-card>
        <div>
          <b-form-textarea
            id="textarea"
            v-model="text"
            placeholder="Enter something..."
            rows="3"
            max-rows="6"
          ></b-form-textarea>
        </div>

        <b-container class="mt-2">
          <b-row>
            <b-col
              cols="12"
              md="4"
            >
            </b-col>
            <b-col
              cols="12"
              md="2"
            >
              <b-button
                @click="speak"
                variant="primary"
                pill
                block
              >发送</b-button>
            </b-col>
            <b-col
              cols="12"
              md="2"
            >
              <b-button
                @click="getContent"
                variant="success"
                pill
                block
              >接收</b-button>
            </b-col>
            <b-col
              cols="12"
              md="4"
            >
            </b-col>
          </b-row>
        </b-container>
      </b-card>
    </b-container>
  </div>
</template>

<script>
// @ is an alias to /src
import { mapActions } from 'vuex';

export default {

  name: 'Home',
  data() {
    return {
      content: '',
      text: '',
      content_id: 1,
    };
  },
  methods: {
    ...mapActions('userModule', { userSpeak: 'speak', userGetContent: 'getContent' }),
    speak() {
      if (this.text === '') {
        console.log('输入为空');
        this.$bvToast.toast('输入不能为空', {
          title: '输入为空',
          variant: 'danger',
          solid: true,
        });
        return;
      }

      this.userSpeak(this.text).then((res) => {
        this.text = '';
        console.log(res.data.data);
      }).catch((err) => {
        console.log(err);
      });
    },
    getContent() {
      this.userGetContent(this.content_id).then((res) => {
        this.content_id = res.data.data.content_id;
        if (!res.data.data.content) {
          return;
        }
        res.data.data.content.forEach((contents) => {
          this.content = `${contents.speaker.name}:${contents.content}\n${this.content}`;
        });
      }).catch((err) => {
        console.log(err);
      });
    },
  },
};
</script>
